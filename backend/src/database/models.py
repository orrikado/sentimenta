from datetime import datetime
from typing import Annotated
from sqlalchemy import ForeignKey, text
from sqlalchemy.orm import Mapped, mapped_column, relationship
from database.db_setup import Base

created_at = Annotated[datetime, mapped_column(server_default=text("TIMEZONE('UTC-4', CURRENT_TIMESTAMP)"))]
updated_at = Annotated[
    datetime,
    mapped_column(
        server_default=text("TIMEZONE('UTC-4', CURRENT_TIMESTAMP)"),
        onupdate=text("TIMEZONE('UTC-4', CURRENT_TIMESTAMP)"),
    ),
]

intpk = Annotated[int, mapped_column(primary_key=True, autoincrement=True)]

class UserOrm(Base):
    __tablename__ = "users"
    uid: Mapped[intpk]
    username: Mapped[str]
    email: Mapped[str]
    password_hash: Mapped[str]
    moods: Mapped[list["MoodOrm"] | None] = relationship(back_populates="user")
    created_at: Mapped[created_at]
    updated_at: Mapped[updated_at]
        
class MoodOrm(Base):
    __tablename__ = "moods"
    uid: Mapped[intpk]
    score: Mapped[int]
    description: Mapped[str]
    date: Mapped[datetime]
    created_at: Mapped[created_at]
    updated_at: Mapped[updated_at]
    user_uid: Mapped[int] = mapped_column(ForeignKey("users.uid"))
    user: Mapped["UserOrm"] = relationship(back_populates="moods")
    