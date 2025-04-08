from pydantic import BaseModel
from datetime import datetime 

class UserSchema(BaseModel):
    uid: int
    username: str
    email: str
    password_hash: str
    created_at: datetime
    updated_at: datetime