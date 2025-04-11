from fastapi import Depends, HTTPException
from fastapi.routing import APIRouter
from database.models import UserOrm
from database.queries.user import get_user
from jwt_token import security

router = APIRouter()

@router.get("/api/user/{id}", dependencies=[Depends(security.access_token_required)])
async def get_user_route(id: int):
    user = await get_user(UserOrm.uid == id)
    if user:
        del user["password_hash"]
        return user
    else:
        raise HTTPException(status_code=404, detail="User not found")
