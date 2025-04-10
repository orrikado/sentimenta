from typing import Optional
from pydantic import BaseModel
from datetime import datetime 
from src.database.schemas.user import UserRelSchema


class MoodSchema(BaseModel):
    uid: int
    score: int
    description: str
    created_at: datetime
    updated_at: datetime

class MoodRelSchema(MoodSchema):
    user: UserRelSchema