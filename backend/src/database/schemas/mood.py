from pydantic import BaseModel
from datetime import datetime
from database.schemas.user import UserRelSchema


class MoodSchema(BaseModel):
    uid: int
    score: int
    description: str
    date: datetime
    created_at: datetime
    updated_at: datetime


class MoodRelSchema(MoodSchema):
    user: UserRelSchema
