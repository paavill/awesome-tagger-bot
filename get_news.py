from selenium import webdriver
from selenium.webdriver.firefox.service import Service
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

import time
import datetime

service = Service('./geckodriver')

options = webdriver.FirefoxOptions()

options.set_preference("dom.webdriver.enabled", False)
options.set_preference("useAutomationExtension", False)
options.add_argument("--disable-dev-shm-usage")
options.add_argument("--disable-gpu")
options.add_argument("--headless")

driver = webdriver.Firefox(service=service, options=options)
driver.get('https://kakoysegodnyaprazdnik.ru')

current_time = datetime.datetime.now().strftime("%Y-%m-%d_%H-%M-%S")

filename = f"./logs/file_{current_time}.txt"
with open(filename, 'w') as file:
    file.write("start...\n")
try:
    element = driver.find_element(By.CLASS_NAME, "mainpage")
    page_html = driver.page_source
    print(page_html)
    with open(filename, 'a') as file:
        file.write("get html in first block\n")
except:
    try:
        submit_button = WebDriverWait(driver, 600).until(
            EC.element_to_be_clickable((By.CSS_SELECTOR, "input[type='submit']"))
        )
        submit_button.click()
        element = WebDriverWait(driver, 600).until(
            EC.presence_of_element_located((By.CLASS_NAME, "mainpage"))
        )
        page_html = driver.page_source
        print(page_html)
        with open(filename, 'a') as file:
            file.write("get html in second block\n")
    except:
        with open(filename, 'a') as file:
            file.write("get html FAILED\n")
finally:
    tabs = driver.window_handles

    for tab in tabs:
        with open(filename, 'a') as file:
            file.write("close tab...\n")
        driver.switch_to.window(tab)
        driver.close()

    driver.quit()