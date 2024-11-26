FROM python:3.12-bookworm
# Set working directory
WORKDIR /app

# RUN apt-get update && apt-get add gcompat

RUN python -m venv /app/venv

# Install Python dependencies
RUN . /app/venv/bin/activate && pip install --upgrade pip && pip install tiktok_captcha_solver 

RUN pip install --upgrade pip && pip install playwright && \
    playwright install --with-deps
RUN . /app/venv/bin/activate && pip install pytest-playwright

# # Cài đặt các trình duyệt cho Playwright (nếu cần)
# RUN . /app/venv/bin/activate && playwright install

# Copy source code
COPY . .

# Run the Python script (adjust command as needed)
# CMD ["/app/venv/bin/python", "scripts/analysis.py"]