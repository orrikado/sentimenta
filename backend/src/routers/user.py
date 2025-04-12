from fastapi import Depends, HTTPException
from fastapi.routing import APIRouter
from database.models import UserOrm
from database.schemas.user import UserUpdateSchema
from database.queries.user import get_user, update_user
from authx import TokenPayload
from jwt_token import security

router = APIRouter(tags=["User"])


@router.get("/api/user/get", dependencies=[Depends(security.access_token_required)])
async def get_user_route(user: TokenPayload = Depends(security.access_token_required)):
    user_id = int(user.sub)
    user = await get_user(UserOrm.uid == user_id)
    if user:
        return user
    else:
        raise HTTPException(status_code=404, detail="User not found")


@router.put("/api/user/update", dependencies=[Depends(security.access_token_required)])
async def update_user_route(
    user_schema: UserUpdateSchema,
    user: TokenPayload = Depends(security.access_token_required),
):
    user_id = int(user.sub)
    try:
        await update_user(UserOrm.uid == user_id, user_schema)
    except Exception as e:
        print(e)
        raise HTTPException(status_code=404, detail="User not found")
