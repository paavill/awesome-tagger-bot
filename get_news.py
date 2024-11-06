from selenium import webdriver
from selenium.webdriver.firefox.service import Service
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
import datetime
import argparse

parser = argparse.ArgumentParser(description="Получение html через Selenium")

parser.add_argument('--url', type=str, help="Адрес сайта, html которого нужно получить", required=True)
parser.add_argument('--proxy', type=str, default="", help="Адрес прокси (если пустой, то прокси не используется)", required=False)
parser.add_argument('--news', type=bool, default=False, help="Указывает ждать ли данных с сайта какой сегодня праздник", required=False)

args = parser.parse_args()
url = args.url
proxy = args.proxy
news = args.news

service = Service('./geckodriver.exe')

options = webdriver.FirefoxOptions()

options.set_preference("dom.webdriver.enabled", False)
options.set_preference("useAutomationExtension", False)
options.add_argument("--disable-dev-shm-usage")
options.add_argument("--disable-gpu")
options.add_argument("--headless")

if proxy != "":
    ip, port = proxy.split(':')
    options.set_preference('network.proxy.type', 1)
    options.set_preference('network.proxy.SOCKS', ip)
    options.set_preference('network.proxy.SOCKS_port', int(port))

driver = webdriver.Firefox(service=service, options=options)
driver.get(url)

if not news:
    print(driver.page_source)
    tabs = driver.window_handles

    for tab in tabs:
        driver.switch_to.window(tab)
        driver.close()

    driver.quit()
    exit()

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
        submit_button = WebDriverWait(driver, 300).until(
            EC.element_to_be_clickable((By.CSS_SELECTOR, "input[type='submit']"))
        )
        submit_button.click()
        element = WebDriverWait(driver, 300).until(
            EC.presence_of_element_located((By.CLASS_NAME, "mainpage"))
        )
        page_html = driver.page_source
        print(page_html)
        with open(filename, 'a') as file:
            file.write("get html in second block\n")
    except:
        with open(filename, 'a') as file:
            file.write("get html FAILED\n")
            file.write(driver.page_source)
finally:
    tabs = driver.window_handles

    for tab in tabs:
        with open(filename, 'a') as file:
            file.write("close tab...\n")
        driver.switch_to.window(tab)
        driver.close()

    driver.quit()