from fastapi.routing import APIRouter
from fastapi import HTTPException, Response
from database.models import User
from database.queries.user import add_user, get_user
from hash import hash_password, verify_password
from src.config import settings
from src.database.schemas.user import UserLoginSchema, UserRegisterSchema
from src.jwt import security

router = APIRouter()

@router.post("/api/auth/login")
async def login(credentials: UserLoginSchema, response: Response):
     user = await get_user(User.email == credentials.email)
     if user and verify_password(credentials.password, user.password_hash):
         jwt_token = security.create_access_token(user.uid)
         response.set_cookie(settings.JWT_COOKIE_NAME, jwt_token)
         return {"access_token": jwt_token}
     else:
         
         raise HTTPException(status_code=401, detail="Incorrect email or password")
     
@router.post("/api/auth/register")
async def register(user_data: UserRegisterSchema, response: Response):
    existing_user = await get_user(User.email == user_data.email)
    if existing_user:
        raise HTTPException(status_code=400, detail="User already exists")
    
    user_data.password = hash_password(user_data.password)
    
    user = await add_user(user=user_data)
    jwt_token = security.create_access_token(user.uid)
    response.set_cookie(settings.JWT_COOKIE_NAME, jwt_token)
    
    return {"access_token": jwt_token}