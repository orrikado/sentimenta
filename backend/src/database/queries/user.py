from database.models import UserOrm
from database.schemas.user import (
    UserAddSchema,
    UserSchema,
    UserUpdateSchema,
    UserWithoutPassword,
)
from database.db_setup import session
from sqlalchemy import select, and_
from sqlalchemy.orm import selectinload


async def get_user(*filters, with_password=True) -> UserSchema:
    async with session() as s:
        stmt = (
            select(UserOrm).where(and_(*filters)).options(selectinload(UserOrm.moods))
        )

        result = await s.execute(stmt)
        user = result.scalar_one_or_none()
        if with_password:
            return (
                UserSchema.model_validate(user, from_attributes=True) if user else None
            )
        else:
            return (
                UserWithoutPassword.model_validate(user, from_attributes=True)
                if user
                else None
            )


async def add_user(user: UserAddSchema) -> UserSchema:
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


async def update_user(user_id: int, user_schema: UserUpdateSchema) -> UserSchema:
    async with session() as s:
        stmt = (
            select(UserOrm)
            .where(UserOrm.uid == user_id)
            .options(selectinload(UserOrm.moods))
        )
        result = await s.execute(stmt)
        user = result.scalar_one()
        update_data = user_schema.model_dump(exclude_unset=True)

        for key, value in update_data.items():
            if key != "uid":
                setattr(user, key, value)

        await s.commit()
        await s.refresh(user)
        user = UserSchema.model_validate(user, from_attributes=True)
        return user
