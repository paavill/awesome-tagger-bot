from selenium import webdriver
from selenium.webdriver.firefox.service import Service
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

import time

service = Service('./geckodriver')

options = webdriver.FirefoxOptions()

options.set_preference("dom.webdriver.enabled", False)
options.set_preference("useAutomationExtension", False)
options.add_argument("--disable-dev-shm-usage")
options.add_argument("--disable-gpu")
options.add_argument("--headless")

driver = webdriver.Firefox(service=service, options=options)

try:
    driver.get('https://kakoysegodnyaprazdnik.ru')
    
    submit_button = WebDriverWait(driver, 600).until(
        EC.element_to_be_clickable((By.CSS_SELECTOR, "input[type='submit']"))
    )
    submit_button.click()
    element = WebDriverWait(driver, 600).until(
        EC.presence_of_element_located((By.CLASS_NAME, "mainpage"))
    )
    page_html = driver.page_source
    print(page_html)
finally:

    driver.quit()