from fastapi import FastAPI
from pydantic import BaseModel
from dotenv import load_dotenv
from functools import lru_cache
from fastapi import Depends, FastAPI
from typing_extensions import Annotated

import os

load_dotenv(".env")

# Get full path to the directory of this file 
# BASEDIR = os.path.abspath(os.path.dirname(__file__))

# @lru_cache
# def get_settings():
#     return config.Settings()


app = FastAPI()



# username, realname, bio, followings, followers, likes, is founder, company name, collaborator_email 
class InfluencerBio(BaseModel):
    description: str
    collaborator_email: str
    is_founder: bool    
    company_name: str 

class Influencer(BaseModel):
    username: str 
    full_name: str 
    bio: InfluencerBio
    avatar: str 
    followers_count: int
    followings_count: int 
    likes_count: int    


@app.get("/")
async def root():
    print("Loi")
    test_var = os.getenv("SCAPTCHA_API_KEY")
    return {"message": "Hello World", "data": test_var, "name": "Loi"}