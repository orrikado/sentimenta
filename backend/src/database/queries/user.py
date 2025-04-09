from src.database.models import User
from src.database.schemas import UserSchema
from src.database.db_setup import session
from sqlalchemy import select
from sqlalchemy.orm import selectinload

async def get_user(uid: int) -> UserSchema:
    async with session() as s: 
        stmt = (
            select(User)
            .where(User.uid == uid)
            .options(selectinload(User.moods))
            )
        result = await s.execute(stmt)
        user = result.scalar_one()
        user = UserSchema.model_validate(user, from_attributes=True)
        return user
    
async def add_user(user: UserSchema) -> UserSchema:
    async with session() as s:
        stmt = (
            select(User)
            .where(User.uid == user.uid)
            .options(selectinload(User.moods))
        )
        result = await s.execute(stmt)
        existing_user = result.scalars().one_or_none()
        if existing_user:
            raise ValueError("User already exists")
        
        user = User(
            username=user.username,
            email=user.email,
            password_hash=user.password_hash,
            moods=[]
        )
        
        s.add(user)
        await s.commit()
        await s.refresh(user)
        user = UserSchema.model_validate(user, from_attributes=True)
        return user

async def delete_user(user_id: int):
    async with session() as s:
        stmt = (
            select(User)
            .where(User.uid == user_id)
        )
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
            select(User)
            .where(User.uid == user.uid)
            .options(selectinload(User.moods))
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