from authx import AuthX, AuthXConfig
from src.config import settings

config = AuthXConfig()
config.JWT_SECRET_KEY = settings.JWT_SECRET
config.JWT_ACCESS_COOKIE_NAME = settings.JWT_COOKIE_NAME
config.JWT_TOKEN_LOCATION = ["cookies"]

security = AuthX(config=config)