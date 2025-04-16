from database.models import MoodOrm
from database.schemas.mood import MoodAddSchema, MoodSchema, MoodUpdateSchema
from database.db_setup import session
from sqlalchemy import select
from sqlalchemy.orm import selectinload


async def get_moods(user_id: int) -> list[MoodSchema]:
    async with session() as s:
        stmt = (
            select(MoodOrm)
            .where(MoodOrm.user_uid == user_id)
            .options(selectinload(MoodOrm.user))
        )
        result = await s.execute(stmt)
        mood = result.scalars().all()
        return [MoodSchema.model_validate(mood, from_attributes=True) for mood in mood]


async def add_mood(user_id: int, mood_schema: MoodAddSchema) -> MoodSchema:
    async with session() as s:
        mood = MoodOrm(
            user_uid=user_id,
            score=mood_schema.score,
            description=mood_schema.description,
            date=mood_schema.date,
            emotions=mood_schema.emotions,
        )

        s.add(mood)
        await s.commit()
        await s.refresh(mood)
        mood = MoodSchema.model_validate(mood, from_attributes=True)
        return mood


async def delete_mood(mood_id: int):
    async with session() as s:
        stmt = select(MoodOrm).where(MoodOrm.uid == mood_id)
        result = await s.execute(stmt)
        mood = result.scalar_one_or_none()
        if mood:
            await s.delete(mood)
            await s.commit()
        else:
            raise ValueError("Mood does not exist")


async def update_mood(user_id: int, mood_schema: MoodUpdateSchema) -> MoodSchema:
    async with session() as s:
        stmt = (
            select(MoodOrm)
            .where(MoodOrm.uid == mood_schema.uid, MoodOrm.user_uid == user_id)
            .options(selectinload(MoodOrm.user))
        )
        result = await s.execute(stmt)
        mood = result.scalar_one()
        update_data = mood_schema.model_dump(exclude_unset=True)

        for key, value in update_data.items():
            if key != "uid":
                setattr(mood, key, value)

        await s.commit()
        await s.refresh(mood)
        mood = MoodSchema.model_validate(mood, from_attributes=True)
        return mood
