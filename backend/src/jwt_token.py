from datetime import timedelta
from authx import AuthX, AuthXConfig
from config import settings

config = AuthXConfig()
config.JWT_SECRET_KEY = settings.JWT_SECRET
config.JWT_ACCESS_COOKIE_NAME = settings.JWT_COOKIE_NAME
config.JWT_TOKEN_LOCATION = ["cookies"]
config.JWT_IMPLICIT_REFRESH_DELTATIME = timedelta(days=30)

security = AuthX(config=config)
