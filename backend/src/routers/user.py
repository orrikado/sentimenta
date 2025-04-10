from fastapi import Depends, HTTPException
from fastapi.routing import APIRouter
from database.models import User
from database.queries.user import get_user
from src.jwt import security

router = APIRouter()

@router.get("/api/user/{id}", dependencies=[Depends(security.access_token_required)])
async def get_user_route(id: int):
    user = await get_user(User.uid == id)
    if user:
        return user
    else:
        raise HTTPException(status_code=404, detail="User not found")
