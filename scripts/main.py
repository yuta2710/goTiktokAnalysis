from fastapi import FastAPI
from pydantic import BaseModel
from dotenv import load_dotenv
from fastapi import FastAPI
from common import parse_formatted_number

from tiktok_captcha_solver import AsyncPlaywrightSolver
from playwright.async_api import async_playwright
from playwright_stealth import StealthConfig, stealth_async
from typing import List


import os
import re 
import time
import uuid

load_dotenv(".env")

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
    followers_count: str
    followings_count: str
    likes_count: str
    


class CommentItem(BaseModel):
    username: str
    content: str 
    full_name: str  
    avatar: str
    # likes_count: str  
    # comments_count: str  
    # saved_count: str 
    # shared_count: str  
    # date_posted: str  
    
    def __hash__(self):
        # Create hashing based on important properties 
        return hash((self.username, self.content, self.avatar))
    
class PostItem(BaseModel):
    username: str 
    content: str
    hashtags: List[str]
    likes_count: str 
    comments_count: str 
    saved_count: str
    shared_count: str
    comments: List[CommentItem]  
    # date_posted: str 
    
    def __hash__(self) -> int:
        return hash((self.content))

class ScrapeProfileRequest(BaseModel):
    url: str

class ScrapePostsRequest(BaseModel):
    url: str 
    number_of_posts: int 
    
class ScrapePostDetailRequest(BaseModel):
    url: str 

@app.post("/scrape")
async def scrapeInfluencerProfile(body: ScrapeProfileRequest):
    api_key = os.getenv("SCAPTCHA_API_KEY")
    
    scraper_limit_posts = 10
    total_views_count = 0
    
    if "DISPLAY" not in os.environ:
        os.environ["DISPLAY"] = ":99"
    
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

        print("\nSolving captcha puzzle in processing ....................")
        await sadcaptcha.solve_captcha_if_present()
        
        print("\nSolved captcha puzzle")

        # Infinite scroll logic
        max_scroll_attempts = 30  # Limit the maximum number of scroll attempts
        scroll_pause_time = 3  # Pause time between scrolls

        seen_view_counts = set()  # Track unique posts already scraped
        total_posts_scraped = 0  # Track total posts scraped
        
        # Username
        username = await page.locator("h1[data-e2e='user-title']").text_content()
        real_name = await page.locator("h2[data-e2e='user-subtitle']").text_content()
        
        avatar = await page.locator("img.css-1zpj2q-ImgAvatar.e1e9er4e1").get_attribute("src")
        
        raw_bio = await page.locator("h2[data-e2e='user-bio']").text_content()
        search_company_name = re.search(r'\S*®', raw_bio)
        
        preprocessed_company_name = ""
        
        if search_company_name:
            preprocessed_company_name = search_company_name.group()

        preprocessed_company_name = preprocessed_company_name.replace("®", "")
        preprocessed_company_name = preprocessed_company_name.lower()

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
            
        for _ in range(max_scroll_attempts):
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
        
        profile_data = Influencer(
            username=username,
            bio=InfluencerBio(
                description=raw_bio,
                collaborator_email=collab_email,
                is_founder=is_founder,
                company_name=preprocessed_company_name,
            ),
            avatar=avatar,
            full_name=real_name,
            followings_count=following_count,
            followers_count=followers_count,
            likes_count=likes_count,
        )

    
    return {
        "success": True,
        "message": "Scrape profile successfully",
        "data": profile_data
    }
     
     
@app.post("/scrape/video/{video_id}")
async def scrapeInfluencerDetailsOfPost(video_id: int, body: ScrapePostDetailRequest):
    api_key = os.getenv("SCAPTCHA_API_KEY")
    tiktok_url = body.url

    if "DISPLAY" not in os.environ:
        os.environ["DISPLAY"] = ":99"
        
    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True)

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
        
        scrape_url = str(tiktok_url + "/" + "video" + "/" + str(video_id))
        # scrape_url = tiktok_url

        # Navigate to the TikTok profile URL
        await page.goto(scrape_url)

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

        total_comments_scraped = 0  # Track total posts scraped
        empty_attempts = 0  # Counter for attempts where no posts are found

        # raw_description = page.locator("span.css-j2a19r-SpanText.efbd9f0").text_content()
        likes_count = 0 
        comments_count = 0 
        saved_count = 0 
        shared_count = 0 
        raw_description = ""
        raw_test_description = ""  
        
        preprocessed_description = ""
        prev_container_text = ""
        
        prev_comments_count = 0 
        retry_count = 0  # To handle retries when no new comments are loaded
        MAX_RETRIES = 6  # Maximum retries before stopping the scroll
        # total_comment_loaded = 0 

        try:
            # Wait for the parent locator containing the spans
            await page.locator("h1[data-e2e='browse-video-desc']").wait_for(timeout=5000)  # Wait for up to 5 seconds
            
            # Fetch the first matching span element
            first_span = page.locator("span.css-j2a19r-SpanText.efbd9f0").first
            await first_span.wait_for(timeout=5000)  # Wait for the specific first element
            
            # Extract the text content of the first span
            preprocessed_description = await first_span.text_content()

        except Exception as e:
            print(f"Error processing description: {e}")

        raw_likes_count = await page.locator("strong[data-e2e='like-count']").text_content()
        raw_comments_count = await page.locator("strong[data-e2e='comment-count']").text_content()
        raw_saved_count = await page.locator("strong[data-e2e='undefined-count']").text_content()
        raw_shared_count = await page.locator("strong[data-e2e='share-count']").text_content()
        
        hashtags = []
        try:
            raw_hashtags = await page.locator("a[data-e2e='search-common-link']").all()
        
            if raw_hashtags is not None:
                for raw_hashtag in raw_hashtags:
                    hashtag_text = await raw_hashtag.locator("strong").text_content()
                    if hashtag_text is not None:
                        hashtags.append(hashtag_text)
        except Exception as e:
            print("\nError scrape hashtags")
        
        print(f"\nInformation of post {video_id}")
        print(raw_description)
        print(preprocessed_description)
        print(hashtags)
        print(raw_comments_count)
        print(raw_saved_count)
        print(raw_shared_count)
        
        limit_comments = int(raw_comments_count.strip())
        
        print(f"\n\nWe've {limit_comments} comments of this post")
        
        comments_data = set()  # Track unique posts already scraped
        comments_list = set()  # Track unique posts already scraped
        processed_subcontents = set()
        
        detect_comment_size = 0 

        # Close the browser after operation (optional)
        # browser.close()

        for _ in range(max_scroll_attempts):
            try:
                if detect_comment_size >= limit_comments:
                    print(f"\nScraping limit of {limit_comments} posts reached. Stopping.")
                    break
                
                current_scroll_height = await page.evaluate("document.body.scrollHeight")
                # await page.wait_for_timeout(8000)
                
                # Scroll to the bottom of the page
                await page.evaluate("window.scrollBy(0, document.body.scrollHeight)")
                await page.wait_for_timeout(5000)  # Allow time for content to load
                
                
                new_scroll_height = await page.evaluate("document.body.scrollHeight")
                
                page.locator("div.css-7whb78-DivCommentListContainer.ezgpko40").wait_for(timeout=5000)
                
                # Locate the container for comments
                container = page.locator("div.css-7whb78-DivCommentListContainer.ezgpko40")
                containerText = await page.locator("div.css-7whb78-DivCommentListContainer.ezgpko40").text_content()
                # Kiểm tra nếu container không thay đổi (không có dữ liệu mới)
                if containerText == prev_container_text:
                    retry_count += 1
                    print(f"No new comments loaded. Retry {retry_count}/{MAX_RETRIES}")
                    if retry_count >= MAX_RETRIES:
                        print("Maximum retries reached. Stopping scroll.")
                        break
                else:
                    retry_count = 0  # Reset retry if there is new data 
                prev_container_text = containerText
                        
                # Ensure comments are loaded
                page.locator(
                    "div.css-13wx63w-DivCommentObjectWrapper.ezgpko42"
                ).wait_for(timeout=10000)
                
                # Extract the comment items
                comment_items = await container.locator(
                    "div.css-13wx63w-DivCommentObjectWrapper.ezgpko42"
                ).all()
                
                total_comments_loaded = len(comment_items)
                prev_comments_count = total_comments_loaded 
                
                print(f"\nLoaded {total_comments_loaded} comments in this attempt.")
                
                await page.evaluate("window.scrollBy(0, document.body.scrollHeight)")
                
                if total_comments_loaded == prev_comments_count:
                    retry_count += 1
                    print(f"No new comments loaded. Retry {retry_count}/{MAX_RETRIES}")
                    if retry_count >= MAX_RETRIES:
                        print("Maximum retries reached. Stopping scroll.")
                        break
                else:
                    retry_count = 0  # Reset retry nếu có comment mới
                    
                prev_comments_count = total_comments_loaded  # Cập nhật số lượng comments đã tải
                new_comments_added = 0

                for cmt_item in comment_items:
                    if total_comments_scraped >= limit_comments:
                        break
                    
                    try:
                        comment_user_order = 'comment-username-1'
                        
                        cmt_username = await cmt_item.locator(f"div[data-e2e='{comment_user_order}'] a.css-22xkqc-StyledLink.er1vbsz0").get_attribute("href")
                        cmt_username = str(cmt_username)
                        cmt_username = cmt_username.replace("/@", "").strip()
                        
                        cmt_fullname = await cmt_item.locator(f"div[data-e2e='{comment_user_order}'] a.css-22xkqc-StyledLink.er1vbsz0 span").text_content()
                        cmt_content = await cmt_item.locator("span[data-e2e='comment-level-1']").text_content()
                        
                        cmt_user_avatar_src = ""
                        cmt_date_posted = ""
                        
                        try:
                            cmt_user_avatar = cmt_item.locator("div.css-vc6h98-DivAvatarWrapper.e1970p9w1 img.css-1zpj2q-ImgAvatar.e1e9er4e1").first
                            cmt_user_avatar_src = await cmt_user_avatar.get_attribute("src")
                        except Exception as e:
                                print("\nError processing avatar of user commentator")
                            
                        unique_key = f"{cmt_username}-{cmt_content}"
                        
                        comment_item_obj = CommentItem(
                            username=cmt_username,
                            content=cmt_content,
                            full_name=cmt_fullname,
                            avatar=cmt_user_avatar_src,
                            # likes_count=raw_likes_count,
                            # comments_count=raw_comments_count,
                            # saved_count=raw_saved_count,
                            # shared_count=raw_shared_count,
                            # date_posted=cmt_date_posted
                        )
                        
                        if unique_key not in comments_data:
                            comments_data.add(unique_key)
                            comments_list.add(comment_item_obj)
                            detect_comment_size += 1 
                            new_comments_added += 1 
                            
                            print(f"\nScraped comment {detect_comment_size}")
                            print(unique_key)
                            
                        # print(cmt_date_posted)
                    except Exception as e:
                        print(f"Error processing next comments: {e}")  # Skip this iteration and continue with the next attempt

                if new_comments_added == 0:
                    retry_count += 1
                    print(f"No new comments loaded. Retry {retry_count}/{MAX_RETRIES}")
                    if retry_count >= MAX_RETRIES:
                        print("Maximum retries reached. Stopping scroll.")
                        break
                else:
                    retry_count = 0  
                    

                if new_scroll_height == current_scroll_height:
                    print("No more content to load. Stopping scroll.")
                    break

            except Exception as e:
                print(f"Error processing comments: {e}")
        await browser.close()
        
        print(f"\n\n\nScraped successfully {len(comments_data)}")
        post_data = PostItem(
            username="louislonghoang",
            content=preprocessed_description,
            hashtags=hashtags,
            likes_count=raw_likes_count,
            comments_count=raw_comments_count,
            saved_count=raw_saved_count,
            shared_count=raw_shared_count,
            comments=comments_list
        )
        
        return {
            "success": True,
            "message": "Scrape profile successfully",
            "count": len(comments_list),
            "data": post_data
        }


# print(
#     bcolors.OKGREEN
#     + f"\nPost {total_posts_scraped} view count: "
#     + bcolors.ENDC,
#     view,
# )
# print(
#     bcolors.OKBLUE + "Content:" + bcolors.ENDC, preprocessed_content
# )
# print(
#     bcolors.OKCYAN + "Hashtags:" + bcolors.ENDC,
#     preprocessed_hashtags,
# )


# comment_item_result = CommentItem(
#     username=cmt_username,
#     full_name=cmt_fullname,
    
# )      

# if view_count not in comments_data:
#     seen_view_counts.add(view_count)
#     total_posts_scraped += 1
#     view = parse_formatted_number(view_count)
                
                
                
                                            # try:
                            #     # Select the nth comment item
                            #     clone_current_cmt_item = page.locator("div.css-7whb78-DivCommentListContainer.ezgpko40") \
                            #         .locator("div.css-13wx63w-DivCommentObjectWrapper.ezgpko42").nth(18)

                            #     # Select the first matching SubContentWrapper within the comment
                            #     subcontent_wrapper = clone_current_cmt_item.locator("div.css-njhskk-DivCommentSubContentWrapper.e1970p9w6").first

                            #     # Gắn UUID cho mỗi subcontent_wrapper
                            #     generated_uuid = str(uuid.uuid4())
                            #     await subcontent_wrapper.evaluate(
                            #         f"element => element.setAttribute('data-uuid', '{generated_uuid}')"
                            #     )

                            #     # Kiểm tra UUID đã gắn
                            #     subcontent_uuid = await subcontent_wrapper.get_attribute("data-uuid")
                            #     print(f"Generated UUID for subcontent_wrapper: {subcontent_uuid}")

                            #     # Check if it's visible
                            #     if await subcontent_wrapper.is_visible():
                            #         print("Subcontent Wrapper is visible.")

                            #         # Interact with this element (e.g., extract text)
                            #         subcontent_text = await subcontent_wrapper.text_content()
                            #         print("Subcontent Text:", subcontent_text)

                            #         # Process the date posted
                            #         cmt_date_posted = subcontent_text.replace("1Reply", "").strip()
                            #         print(f"Date Posted: {cmt_date_posted}")

                            # except Exception as e:
                            #     print(f"Error processing date_posted and likes of comment: {e}")
