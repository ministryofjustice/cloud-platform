#!/usr/bin/env python3
import sys, re, glob
from bs4 import BeautifulSoup
details = glob.glob("rattic/detail*html")
for detail in details:
  with open(detail) as s:
    soup = BeautifulSoup(s, "html5lib")
    t = soup.find("table", attrs={"class" : "table"})
    if t is None:
      sys.exit("failed at {}".format(detail))
    pattern = re.compile("^\s+|\s*,\s*|\s+$")
    for h in t.find_all("th"):
      if h.text == "Owner Group":
        print(h.findNext('td').text)
      if h.text == "Viewer Groups":
        print(pattern.split(h.findNext('td').text.strip()))
      if h.text == "Tags":
        print(pattern.split(h.findNext('td').text.strip()))
