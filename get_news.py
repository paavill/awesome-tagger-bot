from selenium import webdriver
from selenium.webdriver.firefox.service import Service
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
    time.sleep(30)
    driver.get('https://kakoysegodnyaprazdnik.ru')
    page_html = driver.page_source
    print(page_html)
finally:

    driver.quit()