from database.models import UserOrm
from database.schemas.user import UserRegisterSchema, UserSchema
from database.db_setup import session
from sqlalchemy import select, and_
from sqlalchemy.orm import selectinload


async def get_user(*filters) -> UserSchema:
    async with session() as s:
        stmt = (
            select(UserOrm).where(and_(*filters)).options(selectinload(UserOrm.moods))
        )

        result = await s.execute(stmt)
        user = result.scalar_one_or_none()
        return UserSchema.model_validate(user, from_attributes=True) if user else None


async def add_user(user: UserRegisterSchema) -> UserSchema:
    async with session() as s:
        existing_user = await get_user(UserOrm.email == user.email)
        if existing_user:
            raise ValueError("User already exists")

        user = UserOrm(
            username=user.username,
            email=user.email,
            password_hash=user.password,
        )

        s.add(user)
        await s.commit()
        await s.refresh(user)
        user = UserSchema.model_validate(user, from_attributes=True)
        return user


async def delete_user(user_id: int):
    async with session() as s:
        stmt = select(UserOrm).where(UserOrm.uid == user_id)
        result = await s.execute(stmt)
        user = result.scalar_one_or_none()
        if user:
            await s.delete(user)
            await s.commit()
        else:
            raise ValueError("User does not exist")


async def update_user(user: UserSchema, **kwargs) -> UserSchema:
    async with session() as s:
        stmt = (
            select(UserOrm)
            .where(UserOrm.uid == user.uid)
            .options(selectinload(UserOrm.moods))
        )
        result = await s.execute(stmt)
        user = result.scalar_one()
        for key, value in kwargs.items():
            setattr(user, key, value)
        s.add(user)
        await s.commit()
        await s.refresh(user)
        user = UserSchema.model_validate(user, from_attributes=True)
        return user
