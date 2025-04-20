from fastapi.routing import APIRouter
from fastapi import HTTPException, Response
from database.models import UserOrm
from database.queries.user import add_user, get_user
from hash import hash_password, verify_password
from config import settings
from database.schemas.user import UserLoginSchema, UserAddSchema
from jwt_token import security
from utils import validate_email

router = APIRouter(tags=["Auth"])


@router.post("/api/auth/login")
async def login(credentials: UserLoginSchema, response: Response):
    user = await get_user(UserOrm.email == credentials.email)
    if user and verify_password(credentials.password, user.password_hash):
        jwt_token = security.create_access_token(str(user.uid))
        response.set_cookie(settings.JWT_COOKIE_NAME, jwt_token)
        return {"access_token": jwt_token}
    else:
        raise HTTPException(status_code=401, detail="Incorrect email or password")


@router.post("/api/auth/register")
async def register(user_data: UserAddSchema, response: Response):
    existing_user_by_email = await get_user(UserOrm.email == user_data.email)
    if existing_user_by_email:
        raise HTTPException(
            status_code=400, detail="User with this email already exists"
        )

    existing_user_by_username = await get_user(UserOrm.username == user_data.username)
    if existing_user_by_username:
        raise HTTPException(
            status_code=400, detail="User with this username already exists"
        )

    if not validate_email(user_data.email):
        raise HTTPException(status_code=400, detail="Incorrect email")

    user_data.password = hash_password(user_data.password)

    user = await add_user(user=user_data)
    jwt_token = security.create_access_token(str(user.uid))
    response.set_cookie(settings.JWT_COOKIE_NAME, jwt_token)

    return {"access_token": jwt_token}
