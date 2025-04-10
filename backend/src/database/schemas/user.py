from typing import Optional
from pydantic import BaseModel
from datetime import datetime 

class UserSchema(BaseModel):
    uid: int
    username: str
    email: str
    password_hash: str
    created_at: datetime
    updated_at: datetime

class UserRelSchema(UserSchema):
    moods: Optional[list["MoodRelSchema"]] # type: ignore
    
class UserLoginSchema(BaseModel):
    email: str
    password: str

class UserRegisterSchema(BaseModel):
    username: str
    email: str
    password: str