from fastapi import FastAPI
from pydantic import BaseModel
from dotenv import load_dotenv
from functools import lru_cache
from fastapi import Depends, FastAPI
from typing_extensions import Annotated

from common import parse_formatted_number

import os
import re 

from tiktok_captcha_solver import PlaywrightSolver, AsyncPlaywrightSolver
from playwright.sync_api import sync_playwright
from playwright.async_api import async_playwright
from playwright_stealth import stealth_sync, StealthConfig, stealth_async
import time

load_dotenv(".env")

# Get full path to the directory of this file
# BASEDIR = os.path.abspath(os.path.dirname(__file__))

# @lru_cache
# def get_settings():
#     return config.Settings()

class bcolors:
    HEADER = "\033[95m"
    OKBLUE = "\033[94m"
    OKCYAN = "\033[96m"
    OKGREEN = "\033[92m"
    WARNING = "\033[93m"
    FAIL = "\033[91m"
    ENDC = "\033[0m"
    BOLD = "\033[1m"
    UNDERLINE = "\033[4m"

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

class ScrapeProfileRequest(BaseModel):
    url: str

class ScrapePostsRequest(BaseModel):
    url: str 
    number_of_posts: int 

@app.post("/scrape")
async def create_item(body: ScrapeProfileRequest):
    api_key = os.getenv("SCAPTCHA_API_KEY")
    
    scraper_limit_posts = 10
    total_views_count = 0
    
    async with async_playwright() as p:
        # Launch the browser
        browser = await p.chromium.launch(headless=True)
        tiktok_url = body.url 

        page = await browser.new_page()
        await page.add_init_script(
            """
            Object.defineProperty(navigator, 'webdriver', {get: () => undefined});
            """
        )

        # Apply stealth configuration to avoid detection
        config = StealthConfig(
            navigator_languages=False, navigator_vendor=False, navigator_user_agent=False
        )
        await stealth_async(page, config)

        # Navigate to the TikTok profile URL
        await page.goto(tiktok_url)

        # Initialize and use the TikTok CAPTCHA solver
        sadcaptcha = AsyncPlaywrightSolver(
            page,
            api_key,
            mouse_step_size=1,  # Adjust to change mouse speed
            mouse_step_delay_ms=10,  # Adjust to change mouse speed
        )

        await sadcaptcha.solve_captcha_if_present()

        # Infinite scroll logic
        max_scroll_attempts = 30  # Limit the maximum number of scroll attempts
        scroll_pause_time = 3  # Pause time between scrolls

        seen_view_counts = set()  # Track unique posts already scraped
        total_posts_scraped = 0  # Track total posts scraped
        empty_attempts = 0  # Counter for attempts where no posts are found
        
        # Username
        username = await page.locator("h1[data-e2e='user-title']").text_content()
        real_name = await page.locator("h2[data-e2e='user-subtitle']").text_content()
        raw_bio = await page.locator("h2[data-e2e='user-bio']").text_content()
        search_company_name = re.search(r'\S*®', raw_bio)
        
        preprocessed_company_name = ""
        
        if search_company_name:
            preprocessed_company_name = search_company_name.group()

        following_count = await page.locator("strong[data-e2e='following-count']").text_content()
        followers_count = await page.locator("strong[data-e2e='followers-count']").text_content()
        likes_count = await page.locator("strong[data-e2e='likes-count']").text_content()
        
        is_founder = False

        if re.search(r"\b(Founder|CEO)\b", raw_bio, re.IGNORECASE):
            is_founder = True

        email_match = re.search(
            r"\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b", raw_bio
        )

        if email_match:
            collab_email = email_match.group(0)
            
        for attempt in range(max_scroll_attempts):
            if total_posts_scraped >= scraper_limit_posts:
                print(f"\nScraping limit of {scraper_limit_posts} posts reached. Stopping.")
                break

            current_scroll_height = await page.evaluate("document.body.scrollHeight")

            # Scroll to the bottom of the page
            await page.evaluate("window.scrollBy(0, document.body.scrollHeight)")
            await page.wait_for_timeout(3000)  # Allow time for content to load

            new_scroll_height = await page.evaluate("document.body.scrollHeight")

            # Get all posts currently loaded
            post_items = await page.locator('[data-e2e="user-post-item"]').all()

            # Iterate through each post and fetch view counts
            for item in post_items:
                if total_posts_scraped >= scraper_limit_posts:
                    break
                try:
                    # Get view count and ensure uniqueness
                    view_count = await item.locator(
                        'strong[data-e2e="video-views"]'
                    ).text_content()

                    raw_content = await item.locator("div.css-41hm0z img").get_attribute("alt")
                    raw_hashtags = re.findall(r"#\w+", raw_content)

                    preprocessed_hashtags = [tag[1:] for tag in raw_hashtags]

                    preprocessed_content = re.sub(r"#\w+", "", raw_content)
                    preprocessed_content = re.sub(r"#.*", "", raw_content)
                    preprocessed_content = re.sub(
                        r"created by .*", "", preprocessed_content
                    )
                    preprocessed_content = re.sub(r"@.*", "", preprocessed_content)

                    preprocessed_content = preprocessed_content.replace("Ẩn bớt", "")

                    preprocessed_content = preprocessed_content.strip()

                    if view_count not in seen_view_counts:
                        seen_view_counts.add(view_count)
                        total_posts_scraped += 1
                        view = parse_formatted_number(view_count)

                        print(
                            bcolors.OKGREEN
                            + f"\nPost {total_posts_scraped} view count: "
                            + bcolors.ENDC,
                            view,
                        )
                        print(
                            bcolors.OKBLUE + "Content:" + bcolors.ENDC, preprocessed_content
                        )
                        print(
                            bcolors.OKCYAN + "Hashtags:" + bcolors.ENDC,
                            preprocessed_hashtags,
                        )
                        total_views_count += view

                except Exception as e:
                    print(f"Error processing post: {e}")

            # # Scroll to the bottom of the page
            new_scroll_height = await page.evaluate(
                "window.scrollBy(0, document.body.scrollHeight)"
            )
            time.sleep(scroll_pause_time)

            if new_scroll_height == current_scroll_height:
                print("No more content to load. Stopping scroll.")
                break
        
        
        await browser.close()
    
    return {
        "username": username,
        "full_name": real_name,
        "following_count": following_count,
        "followers_count": followers_count,
        "likes_count": likes_count,
        "collab_email": collab_email,
        "is_founder": is_founder,
        "company": preprocessed_company_name
    }
