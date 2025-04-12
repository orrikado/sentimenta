from typing import Optional
from pydantic import BaseModel
from datetime import datetime
from database.schemas.user import UserRelSchema


class MoodSchema(BaseModel):
    uid: int
    score: int
    emotions: str
    description: str | None
    date: datetime
    created_at: datetime
    updated_at: datetime


class MoodRelSchema(MoodSchema):
    user: UserRelSchema


class MoodAddSchema(BaseModel):
    score: int
    emotions: str
    description: str | None = None
    date: datetime


class MoodUpdateSchema(BaseModel):
    uid: int
    score: Optional[int] = None
    emotions: Optional[str] = None
    description: Optional[str] = None
    date: Optional[datetime] = None
