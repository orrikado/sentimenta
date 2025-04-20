from fastapi import Depends, HTTPException
from fastapi.routing import APIRouter
from database.models import UserOrm
from database.schemas.user import UserChangePassSchema, UserUpdateSchema
from database.queries.user import get_user, update_user
from authx import TokenPayload
from jwt_token import security
from hash import hash_password, verify_password

router = APIRouter(tags=["User"])


@router.get("/api/user/get", dependencies=[Depends(security.access_token_required)])
async def get_user_route(user: TokenPayload = Depends(security.access_token_required)):
    user_id = int(user.sub)
    user = await get_user(UserOrm.uid == user_id)
    if user:
        return user
    else:
        raise HTTPException(status_code=404, detail="User not found")


@router.patch(
    "/api/user/update", dependencies=[Depends(security.access_token_required)]
)
async def update_user_route(
    user_schema: UserUpdateSchema,
    user_token: TokenPayload = Depends(security.access_token_required),
):
    user_id = int(user_token.sub)
    user = await get_user(UserOrm.uid == user_id)

    if user_schema.email:
        if not user_schema.password:
            raise HTTPException(status_code=400, detail="Password is required")
        hashed_pass = hash_password(user_schema.password)
        if not verify_password(hashed_pass, user.password_hash):
            raise HTTPException(status_code=400, detail="Incorrect password")
    try:
        await update_user(UserOrm.uid == user_id, user_schema)
    except Exception as e:
        print(e)
        raise HTTPException(status_code=404, detail="User not found")


@router.put(
    "/api/user/update/password", dependencies=[Depends(security.access_token_required)]
)
async def update_user_password_route(
    user_schema: UserChangePassSchema,
    user_token: TokenPayload = Depends(security.access_token_required),
):
    user_id = int(user_token.sub)
    user = await get_user(UserOrm.uid == user_id)

    hashed_pass = hash_password(user_schema.password)
    if not verify_password(hashed_pass, user.password_hash):
        raise HTTPException(status_code=400, detail="Incorrect password")
    else:
        try:
            hashed_new_pass = hash_password(user_schema.new_password)
            await update_user(
                UserOrm.uid == user_id, UserUpdateSchema(password=hashed_new_pass)
            )
        except Exception as e:
            print(e)
            raise HTTPException(status_code=404, detail="User not found")
