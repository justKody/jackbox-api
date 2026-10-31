"""Builds owid_questions.json from Our World in Data charts.

Each entry below names an OWID chart (its URL slug), the entity to read and
the question text. The latest year with data is used as the answer, rounded
to a whole percent. Charts that fail to load or have no data are skipped.

Run: python3 fetch_owid.py
"""

import csv
import io
import json
import sys
import urllib.request
from pathlib import Path

OUT = Path(__file__).with_name("owid_questions.json")
UA = {"User-Agent": "Mozilla/5.0 (jackbox question generator)"}

# (slug, entity, question text)
SPECS = [
    # people
    ("share-of-individuals-using-the-internet", "World", "What percentage of the world's people use the internet?"),
    ("share-of-population-urban", "World", "What percentage of the world's people live in cities?"),
    ("share-of-population-in-extreme-poverty", "World", "What percentage of the world's people live in extreme poverty (under $3 a day)?"),
    ("share-of-the-population-with-access-to-electricity", "World", "What percentage of the world's people have access to electricity?"),
    ("access-to-clean-fuels-and-technologies-for-cooking", "World", "What percentage of the world's people cook with clean fuels like gas or electricity?"),
    ("share-of-adults-who-smoke", "World", "What percentage of the world's adults smoke tobacco?"),
    ("share-of-adults-defined-as-obese", "World", "What percentage of the world's adults are obese?"),
    ("prevalence-of-undernourishment", "World", "What percentage of the world's people are undernourished?"),
    ("literacy-rate-adults", "World", "What percentage of the world's adults can read and write?"),
    ("child-mortality", "World", "What percentage of children worldwide die before their fifth birthday?"),
    ("share-of-individuals-using-the-internet", "India", "What percentage of people in India use the internet?"),
    ("share-of-population-urban", "Japan", "What percentage of people in Japan live in cities?"),
    ("share-of-population-urban", "India", "What percentage of people in India live in cities?"),
    ("share-of-adults-defined-as-obese", "United States", "What percentage of adults in the United States are obese?"),
    ("share-of-adults-who-smoke", "Indonesia", "What percentage of adults in Indonesia smoke tobacco?"),

    # land
    ("forest-area-as-share-of-land-area", "World", "What percentage of the world's land is covered by forest?"),
    ("forest-area-as-share-of-land-area", "Brazil", "What percentage of Brazil's land is covered by forest?"),
    ("forest-area-as-share-of-land-area", "Finland", "What percentage of Finland's land is covered by forest?"),
    ("share-of-land-area-used-for-agriculture", "World", "What percentage of the world's land is used for farming?"),

    # energy
    ("share-electricity-renewables", "World", "What percentage of the world's electricity comes from renewables?"),
    ("share-electricity-coal", "World", "What percentage of the world's electricity comes from coal?"),
    ("share-electricity-nuclear", "World", "What percentage of the world's electricity comes from nuclear power?"),
    ("share-electricity-nuclear", "France", "What percentage of France's electricity comes from nuclear power?"),
    ("share-electricity-renewables", "Norway", "What percentage of Norway's electricity comes from renewables?"),
    ("share-electricity-coal", "China", "What percentage of China's electricity comes from coal?"),
    ("share-electricity-solar", "World", "What percentage of the world's electricity comes from solar power?"),
    ("share-electricity-wind", "Denmark", "What percentage of Denmark's electricity comes from wind power?"),
    ("fossil-fuels-share-energy", "World", "What percentage of the world's energy comes from fossil fuels?"),

    # emissions
    ("annual-share-of-co2-emissions", "China", "China produces what percentage of the world's yearly CO2 emissions?"),
    ("annual-share-of-co2-emissions", "United States", "The United States produces what percentage of the world's yearly CO2 emissions?"),
    ("annual-share-of-co2-emissions", "India", "India produces what percentage of the world's yearly CO2 emissions?"),
    ("annual-share-of-co2-emissions", "Russia", "Russia produces what percentage of the world's yearly CO2 emissions?"),
    ("annual-share-of-co2-emissions", "Japan", "Japan produces what percentage of the world's yearly CO2 emissions?"),

    # more people
    ("share-of-individuals-using-the-internet", "Nigeria", "What percentage of people in Nigeria use the internet?"),
    ("share-of-individuals-using-the-internet", "China", "What percentage of people in China use the internet?"),
    ("share-of-individuals-using-the-internet", "United Kingdom", "What percentage of people in the UK use the internet?"),
    ("share-of-population-urban", "Nigeria", "What percentage of people in Nigeria live in cities?"),
    ("share-of-population-urban", "China", "What percentage of people in China live in cities?"),
    ("share-of-population-urban", "Belgium", "What percentage of people in Belgium live in cities?"),
    ("share-of-population-in-extreme-poverty", "Nigeria", "What percentage of people in Nigeria live in extreme poverty (under $3 a day)?"),
    ("share-of-population-in-extreme-poverty", "India", "What percentage of people in India live in extreme poverty (under $3 a day)?"),
    ("share-of-the-population-with-access-to-electricity", "Nigeria", "What percentage of people in Nigeria have access to electricity?"),
    ("share-of-adults-who-smoke", "United States", "What percentage of adults in the United States smoke tobacco?"),
    ("share-of-adults-who-smoke", "France", "What percentage of adults in France smoke tobacco?"),
    ("share-of-adults-who-smoke", "China", "What percentage of adults in China smoke tobacco?"),
    ("share-of-adults-defined-as-obese", "Japan", "What percentage of adults in Japan are obese?"),
    ("share-of-adults-defined-as-obese", "India", "What percentage of adults in India are obese?"),
    ("share-of-adults-who-are-overweight", "World", "What percentage of the world's adults are overweight or obese?"),
    ("share-of-adults-who-are-overweight", "United States", "What percentage of adults in the United States are overweight or obese?"),
    ("share-of-children-who-are-overweight", "World", "What percentage of children under 5 worldwide are overweight?"),
    ("diabetes-prevalence", "World", "What percentage of the world's adults (ages 20–79) have diabetes?"),
    ("prevalence-of-anemia-in-children", "World", "What percentage of children under 5 worldwide have anemia?"),
    ("literacy-rate-adults", "India", "What percentage of adults in India can read and write?"),
    ("child-mortality", "Nigeria", "What percentage of children in Nigeria die before their fifth birthday?"),
    ("women-survival-to-age-65", "World", "What percentage of women worldwide are expected to live to age 65?"),
    ("men-survival-to-age-65", "World", "What percentage of men worldwide are expected to live to age 65?"),
    ("female-labor-force-participation-rates", "World", "What percentage of the world's women are working or looking for work?"),
    ("proportion-of-women-in-senior-and-middle-management-positions", "United States", "What percentage of senior and middle management jobs in the US are held by women?"),
    ("share-of-births-outside-marriage", "United States", "What percentage of babies in the United States are born to unmarried parents?"),
    ("share-of-births-outside-marriage", "France", "What percentage of babies in France are born to unmarried parents?"),
    ("share-of-births-outside-marriage", "Japan", "What percentage of babies in Japan are born to unmarried parents?"),
    ("how-important-religion-is-in-your-life", "United States", "What percentage of Americans say religion is very important in their life?"),
    ("how-important-religion-is-in-your-life", "Japan", "What percentage of people in Japan say religion is very important in their life?"),
    ("how-important-religion-is-in-your-life", "Nigeria", "What percentage of people in Nigeria say religion is very important in their life?"),
    ("self-reported-trust-attitudes", "United States", "What percentage of Americans agree that \"most people can be trusted\"?"),
    ("self-reported-trust-attitudes", "Norway", "What percentage of Norwegians agree that \"most people can be trusted\"?"),
    ("self-reported-trust-attitudes", "Brazil", "What percentage of Brazilians agree that \"most people can be trusted\"?"),

    # water and schools
    ("proportion-using-safely-managed-drinking-water", "World", "What percentage of the world's people have safe drinking water at home?"),
    ("share-of-the-population-with-access-to-handwashing-facilities", "World", "What percentage of the world's people have soap and water for handwashing at home?"),
    ("agricultural-water-as-a-share-of-total-water-withdrawals", "World", "What percentage of the fresh water humans take from rivers, lakes and the ground is used for farming?"),
    ("primary-schools-with-access-to-internet", "World", "What percentage of the world's primary schools have internet for teaching?"),
    ("percentage-of-primary-schools-with-toilets", "World", "What percentage of the world's primary schools have separate toilets for girls and boys?"),

    # more land and nature
    ("forest-area-as-share-of-land-area", "Japan", "What percentage of Japan's land is covered by forest?"),
    ("forest-area-as-share-of-land-area", "Canada", "What percentage of Canada's land is covered by forest?"),
    ("forest-area-as-share-of-land-area", "Australia", "What percentage of Australia's land is covered by forest?"),
    ("share-of-land-area-used-for-agriculture", "Australia", "What percentage of Australia's land is used for farming?"),
    ("terrestrial-protected-areas", "World", "What percentage of the world's land is in protected areas like national parks?"),
    ("marine-protected-areas", "World", "What percentage of the world's coastal waters are protected?"),
    ("share-of-arable-land-which-is-organic", "World", "What percentage of the world's farmland is organic?"),
    ("share-of-fish-stocks-overexploited", "World", "What percentage of the world's fish stocks are overfished?"),

    # more energy
    ("share-electricity-low-carbon", "World", "What percentage of the world's electricity comes from low-carbon sources (renewables and nuclear)?"),
    ("share-electricity-hydro", "World", "What percentage of the world's electricity comes from hydropower?"),
    ("share-electricity-hydro", "Norway", "What percentage of Norway's electricity comes from hydropower?"),
    ("share-electricity-hydro", "Brazil", "What percentage of Brazil's electricity comes from hydropower?"),
    ("share-electricity-gas", "World", "What percentage of the world's electricity comes from natural gas?"),
    ("share-electricity-gas", "United States", "What percentage of US electricity comes from natural gas?"),
    ("share-electricity-oil", "World", "What percentage of the world's electricity comes from burning oil?"),
    ("share-electricity-nuclear", "United States", "What percentage of US electricity comes from nuclear power?"),
    ("share-electricity-nuclear", "Japan", "What percentage of Japan's electricity comes from nuclear power?"),
    ("share-electricity-renewables", "Germany", "What percentage of Germany's electricity comes from renewables?"),
    ("share-electricity-renewables", "United States", "What percentage of US electricity comes from renewables?"),
    ("share-electricity-coal", "India", "What percentage of India's electricity comes from coal?"),
    ("share-electricity-coal", "Poland", "What percentage of Poland's electricity comes from coal?"),
    ("share-electricity-solar", "Australia", "What percentage of Australia's electricity comes from solar power?"),
    ("share-electricity-wind", "United Kingdom", "What percentage of the UK's electricity comes from wind power?"),
    ("fossil-fuels-share-energy", "United States", "What percentage of US energy comes from fossil fuels?"),
    ("electricity-as-a-share-of-primary-energy", "World", "Electricity makes up what percentage of the world's total energy use?"),
    ("data-centers-share-electricity-demand", "World", "What percentage of the world's electricity is used by data centers?"),
    ("data-centers-share-electricity-demand", "United States", "What percentage of US electricity is used by data centers?"),
    ("electric-car-sales-share", "World", "What percentage of new cars sold worldwide are electric?"),
    ("electric-car-sales-share", "Norway", "What percentage of new cars sold in Norway are electric?"),
    ("electric-car-sales-share", "China", "What percentage of new cars sold in China are electric?"),
    ("electric-car-sales-share", "United States", "What percentage of new cars sold in the United States are electric?"),
    ("military-expenditure-as-a-share-of-government-spending", "United States", "What percentage of US government spending goes to the military?"),
]


def fetch(url):
    req = urllib.request.Request(url, headers=UA)
    with urllib.request.urlopen(req, timeout=60) as r:
        return r.read().decode("utf-8")


def latest(slug, entity, cache):
    if slug not in cache:
        base = f"https://ourworldindata.org/grapher/{slug}"
        reader = csv.DictReader(io.StringIO(fetch(base + ".csv?useColumnShortNames=true")))
        rows = [{k.lower(): v for k, v in row.items()} for row in reader]
        meta = json.loads(fetch(base + ".metadata.json"))
        cache[slug] = (rows, meta)
    rows, meta = cache[slug]

    # first value column after Entity, Code, Year
    col = [c for c in rows[0].keys() if c not in ("entity", "code", "year")][0]
    best = None
    for row in rows:
        if row["entity"] != entity or not row[col]:
            continue
        year = int(row["year"])
        if best is None or year > best[0]:
            best = (year, float(row[col]))
    if best is None:
        raise ValueError(f"no data for {entity}")
    return best, meta["chart"]


def main():
    cache, out = {}, []
    for slug, entity, text in SPECS:
        try:
            (year, value), chart = latest(slug, entity, cache)
        except Exception as e:
            print(f"skip {slug} / {entity}: {e}", file=sys.stderr)
            continue
        answer = round(value)
        if not 0 <= answer <= 100:
            print(f"skip {slug} / {entity}: {value} out of range", file=sys.stderr)
            continue
        out.append({
            "id": f"owid{len(out) + 1}",
            "text": text,
            "answer": answer,
            "exact": round(value, 2),
            "year": year,
            "source": chart.get("citation", ""),
            "url": chart.get("originalChartUrl", f"https://ourworldindata.org/grapher/{slug}"),
        })
        print(f"ok   {entity:15} {year}  {value:6.2f}  {slug}")

    OUT.write_text(json.dumps(out, indent=2, ensure_ascii=False) + "\n")
    print(f"\nwrote {len(out)} questions to {OUT}")


if __name__ == "__main__":
    main()
