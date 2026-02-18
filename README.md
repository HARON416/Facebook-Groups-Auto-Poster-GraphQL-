# Autoposter

Reverse-engineers Facebook's GraphQL API to post product listings (text + images) to your joined groups. It uses browser automation to capture the real API requests, auto-updates the curl commands with fresh cookies and headers, then replays them over HTTP. You just run it and it handles the rest—log in once when prompted, then sit back while it posts across your groups.

## Requirements

- Go 1.25+
- Chrome or Chromium (for browser automation)
- A Facebook account with joined groups

## Setup

### 1. Install dependencies

```bash
go mod tidy
```

### 2. Prepare your items

Create an items directory with one subdirectory per product. Each subdirectory must contain:

- `description.txt` – Post text (one paragraph or multiple lines)
- Image files (jpg, png, etc.) – Up to 6 images per item, randomly selected

Example structure:

```
items/
├── Product1/
│   ├── description.txt
│   ├── image1.jpg
│   └── image2.jpg
├── Product2/
│   ├── description.txt
│   ├── photo1.png
│   └── photo2.png
```

### 3. Configure items path (optional)

By default, Autoposter looks for `./items`. Override with:

- **Command-line flag**: `-items /path/to/your/items`
- **Environment variable**: `ITEMS_PATH=/path/to/your/items`

Example:

```bash
go run . -items /home/user/Downloads/Phones
# or
ITEMS_PATH=/home/user/Downloads/Phones go run .
```

## Login Flow

1. Run the application:

   ```bash
   go run .
   ```

2. A browser window opens to Facebook. **Log in manually** when prompted (email/password, 2FA if enabled).

3. The app waits until it detects a successful login (no login form, no verification screen).

4. It then navigates to your groups and the post composer, capturing the necessary API calls.

5. The browser closes and posting continues in the background via HTTP.

## How It Works

Autoposter reverse engineers Facebook's private GraphQL API (groups fetch, photo upload, post creation) and automatically keeps the captured curl commands up to date with your session.

1. **Capture phase**
   - Opens a visible Chrome window
   - Waits for you to log in to Facebook
   - Navigates to groups page and post composer
   - Intercepts GraphQL and upload requests
   - Extracts curl commands (headers, cookies, endpoints)
   - Closes the browser

2. **Posting phase**
   - Fetches all your joined groups via captured GraphQL
   - Shuffles groups and items for variety
   - For each group (up to 250):
     - Uploads images for the current item
     - Creates a post with description + photos
     - Waits 1 minute between posts
   - Distributes items round-robin across groups

## Output Files

These files are created during the capture phase and contain session data. Do not commit them:

- `graphql_curls.txt` – Groups pagination request
- `graphql_curl.txt` – Create-post mutation
- `uploads_curl.txt` – Image upload request
- `graphql_requests.jsonl` – Raw GraphQL captures
- `uploads.jsonl` – Raw upload captures

`browser_profile/` stores the browser profile (cookies, etc.) per working directory.

## Rate Limiting

- 1 minute delay between posts
- On “We limit how often you…” the app exits
- If you hit limits, wait and run again (the captured curls remain valid until session expiry)

## Configuration Summary

| Option     | Flag / Env   | Default   |
|-----------|--------------|-----------|
| Items path| `-items`     | `./items` |
|           | `ITEMS_PATH` |           |
| Target    | (hardcoded)  | 250 posts |
| Delay     | (hardcoded)  | 1 min     |

## Support & Customizations

- **WhatsApp**: +254718448461
- **Email**: haronkibetrutoh@gmail.com
- **Discord**: kwandapchumba_45230

## License

Private / unlicensed.
