from .auth import router as auth_router
from .user import router as user_router
from .mood import router as mood_router

routers = [auth_router, user_router, mood_router]
