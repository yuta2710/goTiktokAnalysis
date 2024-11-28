from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    app_name: str = "Awesome API"
    admin_email: str
    items_per_user: int = 50
    scaptcha_api_key="91408207650e8ee981b190012199db46" 
    
    model_config = SettingsConfigDict(env_file=".env")
    