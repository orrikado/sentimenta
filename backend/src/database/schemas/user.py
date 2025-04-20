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
    moods: Optional[list["MoodRelSchema"]]  # noqa: F821 # type: ignore


class UserLoginSchema(BaseModel):
    email: str
    password: str


class UserAddSchema(BaseModel):
    username: str
    email: str
    password: str


class UserUpdateSchema(BaseModel):
    username: str | None = None
    email: str | None = None
    password: str | None = None


class UserChangePassSchema(BaseModel):
    password: str
    new_password: str
