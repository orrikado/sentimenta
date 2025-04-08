from pydantic_settings import BaseSettings, SettingsConfigDict

class Settings(BaseSettings):
    
    # Настройки PostgreSQL
    POSTGRES_HOST: str
    POSTGRES_PORT: int = 5432
    POSTGRES_USER: str
    POSTGRES_PASSWORD: str
    POSTGRES_DB: str

    # Строки подключения SQLAlchemy
    SQLALCHEMY_URL: str
    SQLALCHEMY_URL_SYNC: str

    model_config = SettingsConfigDict(env_file=".env", env_file_encoding="utf-8")

    @property
    def SQLALCHEMY_URL(self) -> str:
        # Формирование строки подключения для asyncpg
        return f"postgresql+asyncpg://{self.POSTGRES_USER}:{self.POSTGRES_PASSWORD}@{self.POSTGRES_HOST}:{self.POSTGRES_PORT}/{self.POSTGRES_DB}"

    @property
    def SQLALCHEMY_URL_SYNC(self) -> str:
        # Формирование строки подключения для psycopg2 (синхронный доступ)
        return f"postgresql+psycopg2://{self.POSTGRES_USER}:{self.POSTGRES_PASSWORD}@{self.POSTGRES_HOST}:{self.POSTGRES_PORT}/{self.POSTGRES_DB}"
    
settings = Settings()