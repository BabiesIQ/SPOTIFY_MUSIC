# 🍪 Session Cookie Authentication

> Cookies allow the bot to authenticate requests to external music platforms like YouTube without hitting rate limits or blocks.

---

## ⚠️ Before You Start

- Always use a **throwaway / secondary account** when generating cookies.
- After uploading cookies, **do not log back in** on that account — doing so invalidates the session immediately.
- Cookies expire over time; rotate them periodically.

---

## Step 1 — Export Your Browser Cookies

You need a browser extension to export cookies in **Netscape HTTP format** (`.txt`).

| Browser | Extension | Link |
|---------|-----------|------|
| Chrome | Get cookies.txt LOCALLY | [Install](https://chromewebstore.google.com/detail/get-cookiestxt-clean/ahmnmhfbokciafffnknlekllgcnafnie) |
| Firefox | cookies.txt | [Install](https://addons.mozilla.org/en-US/firefox/addon/cookies-txt/) |

**How to export:**
1. Install the extension above.
2. Open **YouTube.com** in your browser and log in with your secondary account.
3. Click the extension icon → select **"Export"** or **"Get cookies.txt"**.
4. Save the downloaded file somewhere safe.

---

## Step 2 — Host the Cookie File Online

Upload the contents of your `cookies.txt` to a paste service:

| Service | Notes |
|---------|-------|
| [BatBin](https://batbin.me) | ✅ Recommended — no account needed |
| [Pastebin](https://pastebin.com) | Requires free account for persistent pastes |

**Steps:**
1. Open the paste service.
2. Paste the **full content** of `cookies.txt` into the text box.
3. Click **Create / Submit**.
4. Copy the resulting URL (e.g. `https://batbin.me/aBcXyZ`).

---

## Step 3 — Set the Environment Variable

Add the URL to your `.env` file or deployment config:

```env
COOKIES_URL=https://batbin.me/aBcXyZ
```

Multiple cookies (for multiple accounts) are supported — separate by comma:

```env
COOKIES_URL=https://batbin.me/aBcXyZ, https://pastebin.com/xYz123
```

---

## Troubleshooting

| Problem | Fix |
|---------|-----|
| 403 errors still appearing | Cookies may be stale — re-export and re-upload |
| "Session invalid" error | You logged in again on the account — generate fresh cookies |
| Paste URL not working | Make sure the paste is set to **public** / **unlisted** |

---

## Best Practices

- ✅ Rotate cookies every 2–4 weeks
- ✅ Use a private/incognito window when generating cookies
- ✅ Keep a backup paste URL ready
- ❌ Never share your cookie file publicly
