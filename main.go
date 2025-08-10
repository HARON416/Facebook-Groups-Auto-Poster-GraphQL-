package main

import (
	"Facebook-Groups-GraphQL-Auto-Poster/utils"
	"fmt"
	"log"
	"math/rand"
	"time"
)

var fetchGroupsCurl = `curl 'https://web.facebook.com/api/graphql/' \
  -H 'accept: */*' \
  -H 'accept-language: en-US,en;q=0.9' \
  -H 'content-type: application/x-www-form-urlencoded' \
  -b 'sb=XWoHZ-GdXVwlvXZrSaFe7gwz; ps_l=1; ps_n=1; datr=78xqZ8Fc-vdpp5Ii5nGp2A0P; c_user=61560452168137; xs=8%3AeMpQ9UiVlPPogg%3A2%3A1735880128%3A-1%3A-1%3AqSC9VrelgmIhZA%3AAcX2RKQ9XJod2CrR2GMp-u4CJAfKY4Dw6xNT6xwCBhdd; fr=1cDWXPEU6ymLUa9ja.AWceiQdbmfpvQBVg4bPnrF-Xr5D9jdbIJyGakHXNSsgI3VZsdis.Boltxv..AAA.0.0.BolucJ.AWf2q4VIsLoBK35x2GzZuh04KEI; wd=1366x681; presence=C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1754720015009%2C%22v%22%3A1%7D' \
  -H 'origin: https://web.facebook.com' \
  -H 'priority: u=1, i' \
  -H 'referer: https://web.facebook.com/groups/joins/?nav_source=tab&ordering=viewer_added' \
  -H 'sec-ch-prefers-color-scheme: light' \
  -H 'sec-ch-ua: "Not;A=Brand";v="99", "Google Chrome";v="139", "Chromium";v="139"' \
  -H 'sec-ch-ua-full-version-list: "Not;A=Brand";v="99.0.0.0", "Google Chrome";v="139.0.7258.66", "Chromium";v="139.0.7258.66"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-model: ""' \
  -H 'sec-ch-ua-platform: "Linux"' \
  -H 'sec-ch-ua-platform-version: "6.8.0"' \
  -H 'sec-fetch-dest: empty' \
  -H 'sec-fetch-mode: cors' \
  -H 'sec-fetch-site: same-origin' \
  -H 'user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36' \
  -H 'x-asbd-id: 359341' \
  -H 'x-fb-friendly-name: GroupsCometAllJoinedGroupsSectionPaginationQuery' \
  -H 'x-fb-lsd: hafIG2ynVG8MQG0lWxkEiz' \
  --data-raw 'av=61560452168137&__aaid=0&__user=61560452168137&__a=1&__req=1w&__hs=20309.HYP%3Acomet_pkg.2.1...0&dpr=1&__ccg=GOOD&__rev=1025708956&__s=y85dej%3Arv1455%3A4014t0&__hsi=7536465053145952938&__dyn=7xeUjGU5a5Q1ryaxG4Vp41twWwIxu13wFwhUngS3q2ibwNw9G2Saw8i2S1DwUx60GE5O0BU2_CxS320qa321Rwwwqo462mcwfG12wOx62G5Usw9m1YwBgK7o6C0Mo4G17yovwRwlE-U2exi4UaEW2G1jwUBwJK14xm3y11xfxmu3W3iU8o4Wm7-2K1yw9q2-awLyESE2KwwwOg2cwMwhEkxebwHwKG4UrwFg2fwxyo566k1FwgUjz89oeE-3WVU-4FqwIK6E4-mEbUaU2wwgo620XEaUcGy8qxG&__csr=g55fbawR2cchkQ9EtmxJZRcJRdff4W5lOiHkWNnslewyyaOqivtKKAgR9_PaghZTiqj-nZWZAnZLidLQBGQDnEhRuGBWhbOSR_9QKAZ4KhfKnG-utXDFGQ4eh5AF7y4jUSqL-jy4V2ChkiV-GVajyoyqEKAXCnKQECaAVtabFa_zF8CqEKueALVKAF45FpqAKECmidx2l5zGAUzhoVadKiVFFGQcDyKaxa9yonGh92bF1iEKibUlKm7EgG5ud-FpoWp5zbyUoyJ1-5UdVof8yVVK4k364U8E8XxfyUP858iy-EhCz8Cm7oaoW4u4GwBxTwTCxeezpEK2a4pUhxXzoCdx2Unxm8xJ1GqfDwQwiUK4EO0C8gwExK7o8UqwEwFCgyElw8G0HFUxBxVe8wgUlBg4K221Gw-g79DwCU0x-4-5GyQ4ejRgmwxwgqwTyUqw9W5E9U-2S5KewPKVaxbxBa0jmq9w9KcwHx22a684q0y9_bPa0KFFE4pGdF18V83owDxy2y1kyEdU0c180sRw4ew4jwnE0j1K0-E4O15AKcw4fweC3O2O05DEcU1zu01j8w08rUClG1dxq481g81ak0PEfQ0iS0C604sU7skGo3_yk0T89k0PO0Ewb204PK0jJ03f80-a0lt07uAw4Co6a0cvU0-Fv9K1NwRg49w2YE1MQ0oW11xW6iw0AThUC0bzwaGtw4rgmwgo4l0ok04pV2w7lg1YU1HU1E8880EK070E7e0IA0feF0nRg0Za226d0&__hsdp=g4n2LaG64cG1DJaF778Gq2q4E852lgwCgZ9Gy84295lkwIXsOaxqYyA3q4FEmT6-QgOKBeONGihBq1gxx0GzEtFih23RGgO6ghT1z8g8Mi8jMwx2awJbMUxkt2mAyEW92G25F4aAgjCEzKux4lqGcgxDzFryQZBmiXp8WR4vqVK6J4GdudgMKsOx6RhAuZIzF9eFD79chJF4VHlIHlhfjXUPkHqP4J6WGdKXl8ijyXkgvy2YFDiA4d5QqyWgha-hoEPfNXiyF994qcixNpGfm8CyusihAUwXGF64pde9gGql3hqoG4FazFBamtfAyBy-58F2oGazQKleqfy22EyahlpCEDyqKfyF8GQFpmuFpG8GF86zzQ44peh2VoCul48qieK8xa48lGqex6QaAGdixpVoCFWj9yajgJk8zoz8VcI431i2OkmKxKqmArxyfgCca1ng4GcwDghwODxdzkcA8mjBa_ADwh828zU38wu84a1pzA2i2eewHjxG391e7UEUp712fgszV439achA2gGgtgoxb5yoG1ao4mE4R38Qwabwjawj8S55Eai1eu7OwTDw_wbF0sUn8U1j2wDw8K1fw9y2-1dge8eU1oU1gEc83rwd60ri1NwCweu0KGwwAxO1XwoorBDwDwwG6-0xUuw4yw4pwtU5y0yU3QK1ewgU2awrU27wbOaxy1Qw4gwl80Wq0jW0LE1wotxObzo1k85m08axy0RE1J8&__hblp=0rEZ12aBAwko464aof8gwVBwKwYiwEwlp85u1FwMw8e220B8qwxBBwgVU3iwOgSbyUfoW2y2C5XwNxx0zWx279ojz8CfwnoeUuwBg8UC2WE9ebwNyA3C2iUbpEaE2OgG5QbxW2-cwn8Gi58qxG8x24axeUG3OE8o4S2Hxy1kxW222CcxW6UeFoWcgx4Bz8KUy1cyU3Swdu1PwEwpQ2K1hwywxwEzUgwIwGzrxi2npVobo98S0D84y1TyEK8G4K32584l0voO22i0JE3Kw9W1Zwgo88S5RweybwnouUb8eUuxCu0IUbofQt0oA0SUaE7G0D85e11wZwEwbyi0xo4K2m9wCw4eBwl8jwmE2byElwCxnUrBDwi84zCyoW8wQxabyoymuU2exe14g8EcE88rwywYwXwtA4o4q2u1zwjU5a1lwBwExO68S3O1dx-1fwTwMxW0L8G685W5o3nwPAwzwMwOw4Nz89813E2pS8U5u9wbC1yxu4u0GFonwkEdo8E1nUaQ1QzU4K4ouBxV08vwby3ydwvUoxW0JUapVk1eyEeVE9HxGez8doCufwDw&__sjsp=g4n2LaG64cG1DJaF778Gq2q4E852lgwCgZ968wUp4tZl8beTaTFp6L8F74dcX9qqbFdOcTJ49CLir8IAOih_Li8biyCDOoakDBQPrBgh9gBWplz8SVQqEW5qyxGwEixiX8mUOcG7Oqil4AoEKdbFrWAFQER38MyiCF5Kbf5m4UxDxifjKdKS5R4eAdxN7zmfztUAyx6p6y9V6qiULpyFcgnF4VHFCjh9BGbz88krBzrKRisxubJh1-8bEwG4QnhH9F14EFoEP9n7Jafykqcg88yEjCyovxR1a5kaCBwXAG5kmmeAy-bUvgCayEgABgG622Ey5kpGbwMxyp2FoAwahJxB16peiczVUIwuy88UowDzVUa64oGAlohzd1q8wc2Arw9O1bwhUcEaQ4oyi2214w5xwbK1Eg28xG17wMxAs3W5413g591y4Im584Fw9p1Ewb82swzW2Awrw5ug7e0mIE1P84R036Uc80Sa0o20KGwwAxO06-E0fOo&__comet_req=15&fb_dtsg=NAftbcbQfvm30iuXnKtqHKsdH9T9k9kqAL3KHNNS8vHgBJf19vtv6Uw%3A8%3A1735880128&jazoest=25485&lsd=hafIG2ynVG8MQG0lWxkEiz&__spin_r=1025708956&__spin_b=trunk&__spin_t=1754720009&__crn=comet.fbweb.CometGroupsJoinsRoute&fb_api_caller_class=RelayModern&fb_api_req_friendly_name=GroupsCometAllJoinedGroupsSectionPaginationQuery&variables=%7B%22count%22%3A20%2C%22cursor%22%3Anull%2C%22ordering%22%3A%5B%22name%22%5D%2C%22scale%22%3A1%7D&server_timestamps=true&doc_id=9974006939348139'`

var uploadImageCurl = `curl 'https://upload.facebook.com/ajax/react_composer/attachments/photo/upload?av=61560452168137&__aaid=0&__user=61560452168137&__a=1&__req=20&__hs=20309.HYP%3Acomet_pkg.2.1...0&dpr=1&__ccg=MODERATE&__rev=1025709928&__s=i293tt%3Arv1455%3Apvyfye&__hsi=7536487348574657500&__dyn=7xeUjGU5a5Q1ryaxG4Vp41twWwIxu13wFwhUngS3q2ibwNw9G2Saw8i2S1DwUx60GE5O0BU2_CxS320qa321Rwwwqo462mcwfG12wOx62G5Usw9m1YwBgK7o6C0Mo4G17yovwRwlE-U2exi4UaEW2G1jwUBwJK14xm3y11xfxmu2u5Ee88o4Wm7-2K0-obUG2-azqwaW223908O3216xi4UK2K2WEjxK2B08-269wkopg6C13xecwBwWzUfHDzUiBG2OUqwjVqwLwHwa211wo83KwHwOG8xG6E&__csr=g4v6Nc9hkchAbIhbq2tYj5ijh4Y9hYQOJZQCztQBkymy_4YtrOsHGBfRiHGWn948y_46pmy9u_WjGQh8J4RRbEBteGld6KmSWgN_R9aiFaymgxlmp5iLBGQmZeKFpHV642eivgyuXiBy49UVbggm8GjDQnAWhHV9ppCcigXCApl9ABHpUlzkdSWAhfmmmivQ9z9F9Ah94xRoFDy4qniGbF2bDAgyy4DVpogK9U9XK2iaByGG2eGUOKh5CgxmFkaxxe7WGeFx7xuey8pyu7p9ubXml0IyESUnKi4EGm49qxl1qF9oK2xooVFFe48mAyoCudxW6UgBxmUgxeqcCx50GzFUWqu7EK7Hx62aWXxGbz9Qt5Ay8Cu1cxl0HG3-4oSF9o-6WGq1yADG352WwHw-wGxqmbDwEBwRwsU4Ofg4C0wo62Eqxu0oeq17xehwkocEqgr-U7O0Foy1Hx6qmawMAxa5Gw8CGK14wKBwjEBAgqwm85ijwajg4CiE8E5yqt5uhwXwLg2OHt29zAwnFoCKEK0wAEgig_8U7P80L80jFw4uw5pwaWE3Iw1lC06EE0Xe0ge0c5ogK0kSq04_E6aVoK0Ao0VK7riw08Z60j100ge81060i-1exa8O9w7fABwrU1S60QQ09jxq611E6a0lS10y5wEw6Lw-wMw2gS0A87214w4JghwkU0Ed02jEB0kU0jTw-wTw9K042o0E8U0Yp01620UU4K02ny0at80nK4Uay0yK3x0gk2p0l-0fVpUB036Q0nK0reU1BE0ES2-07CUjxC8UkzWp8662l6za6iwMw_o89t06hzE23g7aA04aqA8dgo-1mw&__hsdp=g4b3111a648z8sG8zEuF6Fgwy8ko4H1siF121axAICN0B8htEmT2N2bkG7O9MFIrdkglEg8b49448qG1-QTthG22ibcCxs311D8yEN3EsyEnNb5MEhZgB5AOEAn5iGH8B92PPUFEF2aZmzRSFJk44Qm8WDJ4xqRi4hqrasNBjx6l4ogrIwKm9QAy7JrB8jxCkqP5J7cyip4WV3cTmBaInAeUG8IGFdFi38ynKXx7qgFA4zHcHq9ly5iaagF36CgWAd2yljQhCDGpJnduzDh8WZ6BjpcOgKBzmUFaeh994mE_eOzF44cwGi8gDJ2GUxdxGdrz94bx28EEvWyBmx0x3e7EOEKy22z459m-g92qjwQ8hpj7yAmoEhi8fhCFC6QEgig9kiqy7J5J5DVCFSm4A6_m3OaxOQ1cMO9xly2FyUmBxSt3E9Etg6W2i6UsU4auch9UO1iwKw5UwQwCwloC0AP771OayQbwmQGx-exp3p8D5yU2Bz69wJgtwl452e2G390iE4Wbz821y8ka493iDq6044w7ywm82xa3Oi1Rg-3m0Jo1Uo7q0YE0ye1dwt84q2S0KU8E4a1VwtU6S6Evw8i0lK0PE0iOwaW0fywbm1Ewk82fw2lE1F85-0Uo138&__hblp=0qF4228F0lofUqFwUwjQ3iUeHwxxa16wey0xU4q0wES2ifwj8K2aEdEy7RU98dbwMzEjwPGbx61bxeuewEyp9UK2yQ4GAxm3SA1GxOq2S3Gdx6bg5O4U-12zEmwgo4Ki7E47CiG5KmEfEsxe3GaJ0zz8dU4a1YKqHg8E-fxG7UfXwo89U8EtxWcgZ38uwDwkE21wHwvojxm2e1qwKh8423G7opwZwXDzWwgE4iFE9o668Si9wEx11u586acwjA1twCx-1ixy1hwUwmEOu221yxO0Po5G1Fw-wh83kDzo88-3aUbU5q3mdwjF8C2y7bw8S2e3eqUpw9O3vBzUcUb82kwQwxwNxS0MoC2TwZBwpA3am68bA0X8lK0yE4y2i1zwBxq14xqEW58dE-6Eiz987ebxOfzXwPG5E-u2G2jBwjES12yUF0lfUowooboowd-2qUK1Bxh1q0ZUuwPx61txe1fwg85iexG22dhU-ewKBxGXzryUdo5O3m0EUbE2jxW1WG1wxui3J16U5Foy7obF98y1gCiwqSt3ooy8pwo8lxq1CDwfu9gKegnykA16zAubgcawlU5d0LyEoyrzoK4A0S8yp0qUycwjpUOawzwxBg26hUnyZ0hEmxe7V8ryErxHwc-16BK1Vx26j5gpx61MV9aw&__sjsp=g4b3111a648z8sG8zEuF6Fgwy8ko4H1siF123V4xDsCN0B8WtEmT6RShqRazPOdgFIozl9Vuyk9ll8QhiHCBGkqEAwrJ4JksCcJyfFFFanoCbBWgV0BDhov5SRhFpoWWAxiQ29Zrl7iAhYzowGOkmaVi19233kgx225qfnqCRggjgWcAAxFkx4mCOD9Cle4pkhx1D8bByqF92rGVmm6Vv4INry3dWCjAB8j8ibAwBybaGjqkwlDKXx6uap10-PaSyrgxkyyAaghp3GhRyO6lilh6quFCQBcUih8khFkSegV3QfizAiih5xacx6i4o-i8gDJ2GUCS6ETKcAgK6qa7-EF4FFy8S22cGbBxB1ilLA6O0hy4GkMihpyx58wZ6oBxJa44A2l4yt2QmQmvCqDgpgvm17xOQ1cMO35wyBxSt3E9Etg2jwUwh9UN4Dz81-Ed819NMsyEJ0qWx-2l1P5yU2Bz69wJg7a5i0KwOg1d8dPiDq6044w2TiwYwuA0dkwfa0bCwhEbo2XwywgE7C&__comet_req=15&fb_dtsg=NAfuCk1k9CH1o5HrqxmwwjXApNDJ3jfVpsLsbqxBWj277E2UZrBtVgA%3A8%3A1735880128&jazoest=25510&lsd=2j8Zc1c6DiQIrgGRpx-41s&__spin_r=1025709928&__spin_b=trunk&__spin_t=1754725200&__crn=comet.fbweb.CometGroupDiscussionRoute' \
  -H 'accept: */*' \
  -H 'accept-language: en-US,en;q=0.9' \
  -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary0BBhaiPNQA5HdS3S' \
  -b 'sb=XWoHZ-GdXVwlvXZrSaFe7gwz; ps_l=1; ps_n=1; datr=78xqZ8Fc-vdpp5Ii5nGp2A0P; c_user=61560452168137; wd=1366x681; fr=1lYRCNdKb5Dj3VibZ.AWcihx5ysO5lbZiBIhU3uOM7UqsmVZQZBtOl86uDkBrh3AEleA8.BolvZL..AAA.0.0.BolvZL.AWfO34B692Vc7W5T7pXfgRj45bk; xs=8%3AeMpQ9UiVlPPogg%3A2%3A1735880128%3A-1%3A-1%3AqSC9VrelgmIhZA%3AAcV2K95XU6IN7uvpWUSbL4prLD1OzCCcnpkLiVW0faPr; presence=C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1754725205334%2C%22v%22%3A1%7D' \
  -H 'origin: https://web.facebook.com' \
  -H 'priority: u=1, i' \
  -H 'referer: https://web.facebook.com/' \
  -H 'sec-ch-ua: "Not;A=Brand";v="99", "Google Chrome";v="139", "Chromium";v="139"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "Linux"' \
  -H 'sec-fetch-dest: empty' \
  -H 'sec-fetch-mode: cors' \
  -H 'sec-fetch-site: same-site' \
  -H 'user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36' \
  --data-raw $'------WebKitFormBoundary0BBhaiPNQA5HdS3S\r\nContent-Disposition: form-data; name="source"\r\n\r\n8\r\n------WebKitFormBoundary0BBhaiPNQA5HdS3S\r\nContent-Disposition: form-data; name="profile_id"\r\n\r\n61560452168137\r\n------WebKitFormBoundary0BBhaiPNQA5HdS3S\r\nContent-Disposition: form-data; name="waterfallxapp"\r\n\r\ncomet\r\n------WebKitFormBoundary0BBhaiPNQA5HdS3S\r\nContent-Disposition: form-data; name="farr"; filename="451833410_122123902592349689_2546841386593894967_n.jpg"\r\nContent-Type: image/jpeg\r\n\r\n\r\n------WebKitFormBoundary0BBhaiPNQA5HdS3S\r\nContent-Disposition: form-data; name="upload_id"\r\n\r\njsc_c_l\r\n------WebKitFormBoundary0BBhaiPNQA5HdS3S--\r\n'`

var createPostCurl = `curl 'https://web.facebook.com/api/graphql/' \
  -H 'accept: */*' \
  -H 'accept-language: en-US,en;q=0.9' \
  -H 'content-type: application/x-www-form-urlencoded' \
  -b 'sb=XWoHZ-GdXVwlvXZrSaFe7gwz; ps_l=1; ps_n=1; datr=78xqZ8Fc-vdpp5Ii5nGp2A0P; c_user=61560452168137; wd=1366x681; presence=C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1754730356835%2C%22v%22%3A1%7D; fr=14ozVsshNDEteBKM5.AWfIQN9wRFTsxT-ytgrQmJAtpYsqDytY9gerjtODwrbEoAkFEKw.BolxMO..AAA.0.0.BolxMO.AWdJbLf1AVdSnM-BJXJlo4Y0otU; xs=8%3AeMpQ9UiVlPPogg%3A2%3A1735880128%3A-1%3A-1%3AqSC9VrelgmIhZA%3AAcU7FQ3GK77T_SV2EN1DHjeyC_zPe1yifXly-rJXA5EN' \
  -H 'origin: https://web.facebook.com' \
  -H 'priority: u=1, i' \
  -H 'referer: https://web.facebook.com/groups/818100786185994' \
  -H 'sec-ch-prefers-color-scheme: light' \
  -H 'sec-ch-ua: "Not;A=Brand";v="99", "Google Chrome";v="139", "Chromium";v="139"' \
  -H 'sec-ch-ua-full-version-list: "Not;A=Brand";v="99.0.0.0", "Google Chrome";v="139.0.7258.66", "Chromium";v="139.0.7258.66"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-model: ""' \
  -H 'sec-ch-ua-platform: "Linux"' \
  -H 'sec-ch-ua-platform-version: "6.8.0"' \
  -H 'sec-fetch-dest: empty' \
  -H 'sec-fetch-mode: cors' \
  -H 'sec-fetch-site: same-origin' \
  -H 'user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36' \
  -H 'x-asbd-id: 359341' \
  -H 'x-fb-friendly-name: ComposerStoryCreateMutation' \
  -H 'x-fb-lsd: 2j8Zc1c6DiQIrgGRpx-41s' \
  --data-raw 'av=61560452168137&__aaid=0&__user=61560452168137&__a=1&__req=47&__hs=20309.HYP%3Acomet_pkg.2.1...0&dpr=1&__ccg=MODERATE&__rev=1025709928&__s=g4p73a%3Arv1455%3Apvyfye&__hsi=7536487348574657500&__dyn=7xeUjGU5a5Q1ryaxG4Vp41twWwIxu13wFwhUngS3q2ibwNw9G2Saw8i2S1DwUx60GE5O0BU2_CxS320qa321Rwwwqo462mcwfG12wOx62G5Usw9m1YwBgK7o6C0Mo4G17yovwRwlE-U2exi4UaEW2G1jwUBwJK14xm3y11xfxmu2u5Ee88o4Wm7-2K0-obUG2-azqwaW223908O3216xi4UK2K2WEjxK2B08-269wkopg6C13xecwBwWzUfHDzUiBG2OUqwjVqwLwHwa211wo83KwHwOG8xG6E&__csr=g4v6Nc9hkchAbIhbq2tYj5ijh4Y9hYQOJZQCztQBkymy_4YtrOsHGBfRiHGWn948y_46pmy9u_WjGQh8J4RRbEBteGld6KmSWgN_R9aiFaymgxlmp5iLBGQmZeKFpHV642eivgyuXiBy49UVbggm8GjDQnAWhHV9ppCcigXCApl9ABHpUlzkdSWAhfmmmivQ9z9F9Ah94xRoFDy4qniGbF2bDAgyy4DVpogK9U9XK2iaByGG2eGUOKh5CgxmFkaxxe7WGeFx7xuey8pyu7p9ubXml0IyESUnKi4EGm49qxl1qF9oK2xooVFFe48mAyoCudxW6UgBxmUgxeqcCx50GzFUWqu7EK7Hx62aWXxGbz9Qt5Ay8Cu1cxl0HG3-4oSF9o-6WGq1yADG352WwHw-wGxqmbDwEBwRwsU4Ofg4C0wo62Eqxu0oeq17xehwkocEqgr-U7O0Foy1Hx6qmawMAxa5Gw8CGK14wKBwjEBAgqwm85ijwajg4CiE8E5yqt5uhwXwLg2OHt29zAwnFoCKEK0wAEgig_8U7P80L80jFw4uw5pwaWE3Iw1lC06EE0Xe0ge0c5ogK0kSq04_E6aVoK0Ao0VK7riw08Z60j100ge81060i-1exa8O9w7fABwrU1S60QQ09jxq611E6a0lS10y5wEw6Lw-wMw2gS0A87214w4JghwkU0Ed02jEB0kU0jTw-wTw9K042o0E8U0Yp01620UU4K02ny0at80nK4Uay0yK3x0gk2p0l-0fVpUB036Q0nK0reU1BE0ES2-07CUjxC8UkzWp8662l6za6iwMw_o89t06hzE23g7aA04aqA8dgo-1mw&__hsdp=g4b3111a648z8sG8zEuF6Fgwy8ko4H1siF121axAICN0B8htEmT2N2bkG7O9MFIrdkglEg8b49448qG1-QTthG22ibcCxs311D8yEN3EsyEnNb5MEhZgB5AOEAn5iGH8B92PPUFEF2aZmzRSFJk44Qm8WDJ4xqRi4hqrasNBjx6l4ogrIwKm9QAy7JrB8jxCkqP5J7cyip4WV3cTmBaInAeUG8IGFdFi38ynKXx7qgFA4zHcHq9ly5iaagF36CgWAd2yljQhCDGpJnduzDh8WZ6BjpcOgKBzmUFaeh994mE_eOzF44cwGi8gDJ2GUxdxGdrz94bx28EEvWyBmx0x3e7EOEKy22z459m-g92qjwQ8hpj7yAmoEhi8fhCFC6QEgig9kiqy7J5J5DVCFSm4A6_m3OaxOQ1cMO9xly2FyUmBxSt3E9Etg6W2i6UsU4auch9UO1iwKw5UwQwCwloC0AP771OayQbwmQGx-exp3p8D5yU2Bz69wJgtwl452e2G390iE4Wbz821y8ka493iDq6044w7ywm82xa3Oi1Rg-3m0Jo1Uo7q0YE0ye1dwt84q2S0KU8E4a1VwtU6S6Evw8i0lK0PE0iOwaW0fywbm1Ewk82fw2lE1F85-0Uo138&__hblp=0qF4228F0lofUqFwUwjQ3iUeHwxxa16wey0xU4q0wES2ifwj8K2aEdEy7RU98dbwMzEjwPGbx61bxeuewEyp9UK2yQ4GAxm3SA1GxOq2S3Gdx6bg5O4U-12zEmwgo4Ki7E47CiG5KmEfEsxe3GaJ0zz8dU4a1YKqHg8E-fxG7UfXwo89U8EtxWcgZ38uwDwkE21wHwvojxm2e1qwKh8423G7opwZwXDzWwgE4iFE9o668Si9wEx11u586acwjA1twCx-1ixy1hwUwmEOu221yxO0Po5G1Fw-wh83kDzo88-3aUbU5q3mdwjF8C2y7bw8S2e3eqUpw9O3vBzUcUb82kwQwxwNxS0MoC2TwZBwpA3am68bA0X8lK0yE4y2i1zwBxq14xqEW58dE-6Eiz987ebxOfzXwPG5E-u2G2jBwjES12yUF0lfUowooboowd-2qUK1Bxh1q0ZUuwPx61txe1fwg85iexG22dhU-ewKBxGXzryUdo5O3m0EUbE2jxW1WG1wxui3J16U5Foy7obF98y1gCiwqSt3ooy8pwo8lxq1CDwfu9gKegnykA16zAubgcawlU5d0LyEoyrzoK4A0S8yp0qUycwjpUOawzwxBg26hUnyZ0hEmxe7V8ryErxHwc-16BK1Vx26j5gpx61MV9aw&__sjsp=g4b3111a648z8sG8zEuF6Fgwy8ko4H1siF123V4xDsCN0B8WtEmT6RShqRazPOdgFIozl9Vuyk9ll8QhiHCBGkqEAwrJ4JksCcJyfFFFanoCbBWgV0BDhov5SRhFpoWWAxiQ29Zrl7iAhYzowGOkmaVi19233kgx225qfnqCRggjgWcAAxFkx4mCOD9Cle4pkhx1D8bByqF92rGVmm6Vv4INry3dWCjAB8j8ibAwBybaGjqkwlDKXx6uap10-PaSyrgxkyyAaghp3GhRyO6lilh6quFCQBcUih8khFkSegV3QfizAiih5xacx6i4o-i8gDJ2GUCS6ETKcAgK6qa7-EF4FFy8S22cGbBxB1ilLA6O0hy4GkMihpyx58wZ6oBxJa44A2l4yt2QmQmvCqDgpgvm17xOQ1cMO35wyBxSt3E9Etg2jwUwh9UN4Dz81-Ed819NMsyEJ0qWx-2l1P5yU2Bz69wJg7a5i0KwOg1d8dPiDq6044w2TiwYwuA0dkwfa0bCwhEbo2XwywgE7C&__comet_req=15&fb_dtsg=NAfuCk1k9CH1o5HrqxmwwjXApNDJ3jfVpsLsbqxBWj277E2UZrBtVgA%3A8%3A1735880128&jazoest=25510&lsd=2j8Zc1c6DiQIrgGRpx-41s&__spin_r=1025709928&__spin_b=trunk&__spin_t=1754725200&__crn=comet.fbweb.CometGroupDiscussionRoute&fb_api_caller_class=RelayModern&fb_api_req_friendly_name=ComposerStoryCreateMutation&variables=%7B%22input%22%3A%7B%22composer_entry_point%22%3A%22hosted_inline_composer%22%2C%22composer_source_surface%22%3A%22group%22%2C%22composer_type%22%3A%22group%22%2C%22logging%22%3A%7B%22composer_session_id%22%3A%22d48e55fc-b3b3-4038-944f-b9b40c0a3a7b%22%7D%2C%22source%22%3A%22WWW%22%2C%22message%22%3A%7B%22ranges%22%3A%5B%5D%2C%22text%22%3A%22LIPA%20POLE%20POLE%20Samsung%20Galaxy%20Note%2010%20Plus%20%20Renewed%204G%20%7C%20single%20SIM%20%7C%204300%20mAh%20%7C%206.67%5C%22%20%7C%2012GB%20RAM%20%7C%20256GB%20ROM%20Deposit%20Ksh.%209%2C69912%20month%20Weekly%20Payment%20Ksh.%201360Call%2FWhatsApp%200718448461%22%7D%2C%22with_tags_ids%22%3Anull%2C%22inline_activities%22%3A%5B%5D%2C%22text_format_preset_id%22%3A%220%22%2C%22group_flair%22%3A%7B%22flair_id%22%3Anull%7D%2C%22attachments%22%3A%5B%7B%22photo%22%3A%7B%22id%22%3A%22122200904204348405%22%7D%7D%2C%7B%22photo%22%3A%7B%22id%22%3A%22122200904186348405%22%7D%7D%2C%7B%22photo%22%3A%7B%22id%22%3A%22122200904162348405%22%7D%7D%5D%2C%22composed_text%22%3A%7B%22block_data%22%3A%5B%22%7B%7D%22%2C%22%7B%7D%22%2C%22%7B%7D%22%2C%22%7B%7D%22%2C%22%7B%7D%22%2C%22%7B%7D%22%5D%2C%22block_depths%22%3A%5B0%2C0%2C0%2C0%2C0%2C0%5D%2C%22block_types%22%3A%5B0%2C0%2C0%2C0%2C0%2C0%5D%2C%22blocks%22%3A%5B%22LIPA%20POLE%20POLE%20%22%2C%22Samsung%20Galaxy%20Note%2010%20Plus%20%20%22%2C%22Renewed%204G%20%7C%20single%20SIM%20%7C%204300%20mAh%20%7C%206.67%5C%22%20%7C%2012GB%20RAM%20%7C%20256GB%20ROM%20%22%2C%22Deposit%20Ksh.%209%2C699%22%2C%2212%20month%20Weekly%20Payment%20Ksh.%201360%22%2C%22Call%2FWhatsApp%200718448461%22%5D%2C%22entities%22%3A%5B%22%5B%5D%22%2C%22%5B%5D%22%2C%22%5B%5D%22%2C%22%5B%5D%22%2C%22%5B%5D%22%2C%22%5B%5D%22%5D%2C%22entity_map%22%3A%22%7B%7D%22%2C%22inline_styles%22%3A%5B%22%5B%5D%22%2C%22%5B%5D%22%2C%22%5B%5D%22%2C%22%5B%5D%22%2C%22%5B%5D%22%2C%22%5B%5D%22%5D%7D%2C%22navigation_data%22%3A%7B%22attribution_id_v2%22%3A%22CometGroupDiscussionRoot.react%2Ccomet.group%2Cunexpected%2C1754725204848%2C201917%2C2361831622%2C%2C%3BGroupsCometJoinsRoot.react%2Ccomet.groups.joins%2Cvia_cold_start%2C1754725202145%2C768689%2C%2C%2C%22%7D%2C%22tracking%22%3A%5Bnull%5D%2C%22event_share_metadata%22%3A%7B%22surface%22%3A%22newsfeed%22%7D%2C%22audience%22%3A%7B%22to_id%22%3A%22818100786185994%22%7D%2C%22actor_id%22%3A%2261560452168137%22%2C%22client_mutation_id%22%3A%221%22%7D%2C%22feedLocation%22%3A%22GROUP%22%2C%22feedbackSource%22%3A0%2C%22focusCommentID%22%3Anull%2C%22gridMediaWidth%22%3Anull%2C%22groupID%22%3Anull%2C%22scale%22%3A1%2C%22privacySelectorRenderLocation%22%3A%22COMET_STREAM%22%2C%22checkPhotosToReelsUpsellEligibility%22%3Afalse%2C%22renderLocation%22%3A%22group%22%2C%22useDefaultActor%22%3Afalse%2C%22inviteShortLinkKey%22%3Anull%2C%22isFeed%22%3Afalse%2C%22isFundraiser%22%3Afalse%2C%22isFunFactPost%22%3Afalse%2C%22isGroup%22%3Atrue%2C%22isEvent%22%3Afalse%2C%22isTimeline%22%3Afalse%2C%22isSocialLearning%22%3Afalse%2C%22isPageNewsFeed%22%3Afalse%2C%22isProfileReviews%22%3Afalse%2C%22isWorkSharedDraft%22%3Afalse%2C%22hashtag%22%3Anull%2C%22canUserManageOffers%22%3Afalse%2C%22__relay_internal__pv__CometUFIShareActionMigrationrelayprovider%22%3Atrue%2C%22__relay_internal__pv__GHLShouldChangeSponsoredDataFieldNamerelayprovider%22%3Atrue%2C%22__relay_internal__pv__GHLShouldChangeAdIdFieldNamerelayprovider%22%3Atrue%2C%22__relay_internal__pv__CometUFI_dedicated_comment_routable_dialog_gkrelayprovider%22%3Afalse%2C%22__relay_internal__pv__IsWorkUserrelayprovider%22%3Afalse%2C%22__relay_internal__pv__CometUFIReactionsEnableShortNamerelayprovider%22%3Afalse%2C%22__relay_internal__pv__FBReels_deprecate_short_form_video_context_gkrelayprovider%22%3Atrue%2C%22__relay_internal__pv__FeedDeepDiveTopicPillThreadViewEnabledrelayprovider%22%3Afalse%2C%22__relay_internal__pv__FBReels_enable_view_dubbed_audio_type_gkrelayprovider%22%3Afalse%2C%22__relay_internal__pv__CometImmersivePhotoCanUserDisable3DMotionrelayprovider%22%3Afalse%2C%22__relay_internal__pv__WorkCometIsEmployeeGKProviderrelayprovider%22%3Afalse%2C%22__relay_internal__pv__IsMergQAPollsrelayprovider%22%3Afalse%2C%22__relay_internal__pv__FBReelsMediaFooter_comet_enable_reels_ads_gkrelayprovider%22%3Atrue%2C%22__relay_internal__pv__StoriesArmadilloReplyEnabledrelayprovider%22%3Atrue%2C%22__relay_internal__pv__FBReelsIFUTileContent_reelsIFUPlayOnHoverrelayprovider%22%3Atrue%2C%22__relay_internal__pv__GHLShouldChangeSponsoredAuctionDistanceFieldNamerelayprovider%22%3Atrue%7D&server_timestamps=true&doc_id=24752350781035639'`

func main() {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	phonesPath := "/home/kwandapchumba/Pictures/SIMU"

	err := utils.UpdateFacebookGroupsFromCurl(fetchGroupsCurl, "utils/fetchgroups.go")
	if err != nil {
		log.Fatal(err)
	}

	groups, err := utils.FetchGroups(r)
	if err != nil {
		log.Fatal(err)
	}

	r.Shuffle(len(groups), func(i, j int) {
		groups[i], groups[j] = groups[j], groups[i]
	})

	fmt.Println("Total groups: ", len(groups))

	groups = groups[:200]

	totalPosts := 0

	for _, group := range groups {
		groupName := group.Name
		groupID := group.ID
		groupURL := group.URL

		fmt.Printf("Group ID: %s\nGroup Name: %s\nGroup URL: %s\n", groupID, groupName, groupURL)

		phones, err := utils.ExtractPhones(phonesPath)
		if err != nil {
			log.Fatal(err)
		}

		r.Shuffle(len(phones), func(i, j int) {
			phones[i], phones[j] = phones[j], phones[i]
		})

		phone := phones[0]

		photoIDs := []string{}

		err = utils.UpdateImageUploadFunctionFromCurl(uploadImageCurl, "utils/uploadimage.go")
		if err != nil {
			log.Fatal(err)
		}

		for _, imagePath := range phone.ImagePaths {
			imageID, err := utils.UploadImage(imagePath)
			if err != nil {
				log.Fatal(err)
			}
			photoIDs = append(photoIDs, imageID)
		}

		fmt.Printf("Photo IDs: %v\n", photoIDs)

		post := utils.FacebookPost{
			MessageText: phone.Description,
			PhotoIDs:    photoIDs,
			GroupID:     group.ID,
		}

		err = utils.UpdateCreateGroupPostFunctionFromCurl(createPostCurl, "utils/creategrouppost.go")
		if err != nil {
			log.Fatal(err)
		}

		post, err = utils.CreateGroupPost(post)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v created\n Group Name: %s\n Group URL: %s\n", post, groupName, groupURL)

		totalPosts++

		fmt.Println("Total posts created so far: ", totalPosts)

		duration := utils.ReturnRandomNumber(r, 2, 3)

		time.Sleep(time.Duration(duration) * time.Minute)

		fmt.Printf("Sleeping for %f minutes\n", duration)
	}
}
