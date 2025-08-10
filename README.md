# 📱 Facebook Groups GraphQL Auto Poster

A Go application that posts to Facebook Groups using GraphQL. It reads your content from a local directory, uploads photos, and creates a post per group. The app now uses embedded “Copy as cURL” requests inside `main.go` to auto-update the HTTP request code in `utils/*` before running.

## ✨ Features

- **🔄 Automatic Group Discovery**: Fetches your joined groups via GraphQL with pagination
- **📸 Image Upload**: Uploads multiple photos and collects their IDs
- **📝 Content Loader**: Reads descriptions and images from a directory structure
- **⏱️ Scheduling**: Random delay between posts (2–3 minutes)
- **🎲 Randomization**: Shuffles groups and items for natural behavior
- **🧠 Self-updating Requests**: Parses your fresh browser cURL and updates code in `utils/*` automatically

## 🏗️ Current Architecture

```
POST_WITH_GRAPHQL/
├── main.go                 # Entrypoint (paste fresh cURL strings here)
├── go.mod
├── go.sum
├── makefile
└── utils/
    ├── fetchgroups.go                         # Fetch joined groups (GraphQL)
    ├── creategrouppost.go                     # Create group post (GraphQL)
    ├── uploadimage.go                         # Upload image endpoint
    ├── extractphones.go                       # Content loader from directories
    ├── buildattachments.go                    # Build attachments JSON from photo IDs
    ├── returnrandomnumber.go                  # Random delay generator
    ├── updatefetchgroupscookies.go            # Update fetchgroups.go from cURL
    ├── updateuploadphotocookies.go            # Update uploadimage.go from cURL
    └── updatecreategrouppostcookies.go        # Update creategrouppost.go from cURL
```

Backups are created automatically the first time each updater runs: `*.go.backup`.

## 🚀 Getting Started

### Prerequisites

- Go 1.24.5+
- A Facebook account with joined groups
- A local content directory

### Install

```bash
git clone <repository-url>
cd POST_WITH_GRAPHQL
go mod tidy
```

### Configure

1. Set your content path in `main.go`:

```go
phonesPath := "/absolute/path/to/CONTENT_ROOT"
```

2. Paste fresh “Copy as cURL” strings into `main.go` variables:

- `fetchGroupsCurl` (GraphQL joined groups request)
- `uploadImageCurl` (upload endpoint with multipart form)
- `createPostCurl` (ComposerStoryCreateMutation GraphQL)

How to obtain:

- Open Facebook in a logged-in browser
- Open DevTools → Network
- Perform the action (open groups page, start a post with photos, submit a post)
- Right‑click the relevant request → Copy → Copy as cURL
- Paste the entire cURL into the corresponding variable in `main.go`

No `cookies/` directory is used anymore.

### Run

```bash
make run
# or
go run main.go
```

## 📖 How It Works (current behavior)

- The app updates code in `utils/*` at runtime using the three cURL strings:
  - `UpdateFacebookGroupsFromCurl` → writes fresh params/headers/cookies into `utils/fetchgroups.go`
  - `UpdateImageUploadFunctionFromCurl` → writes fresh params/headers/cookies into `utils/uploadimage.go`
  - `UpdateCreateGroupPostFunctionFromCurl` → writes fresh params/headers/cookies into `utils/creategrouppost.go`
- Groups are fetched, shuffled, and then limited: `groups = groups[:200]`.
- For each group, the app:
  - Loads phones/items from your content directory
  - Shuffles items and picks one item (`phones[0]`) for this group
  - Uploads all images in that item and collects photo IDs
  - Creates a post with the item’s description and the uploaded photo IDs
  - Sleeps a random 2–3 minutes before the next group

## 📂 Content Structure

```
CONTENT_ROOT/
├── item1/
│   ├── details.txt            # contains:  "description: ..."
│   ├── image1.jpg
│   └── image2.jpg
├── item2/
│   ├── details.txt
│   └── image1.jpg
└── ...
```

`details.txt` format:

```
description: Your text...feature one...feature two...feature three
```

The app uppercases the description, splits by `...`, and formats each part with a leading ✅ and blank lines between.

## 🔧 Technical Notes

- **GraphQL**: Uses `GroupsCometAllJoinedGroupsSectionPaginationQuery` to list joined groups; uses Composer mutation for posting
- **Uploads**: Uses the Facebook upload endpoint to obtain real `photoID`s and validates image URLs in the response
- **Self-updating requests**: Updaters replace request bodies, headers, cookies, LSD/FB-DTSG tokens, etc., then write `*.backup` once
- **Rate/safety**: Sleeps for a random 2–3 minutes between posts

## 🛡️ Tips & Troubleshooting

- Use very fresh cURL copies; Facebook tokens expire quickly
- Ensure the cURL you paste corresponds to the same account/session
- If you get 403/redirects or empty data, refresh cURL and try again
- Backups are saved next to the updated files as `*.go.backup`

## 📋 Requirements

- Go 1.24.5+
- Dependency: `github.com/klauspost/compress v1.18.0`
- Facebook account and fresh browser session
- Local content directory

## 🛠️ Development

```bash
make run     # run the app
make tidy    # go mod tidy
go build -o facebook-poster main.go
```

Key modules:

- `main.go` orchestrates fetching → uploading → posting
- `utils/fetchgroups.go` GraphQL groups listing
- `utils/uploadimage.go` image upload
- `utils/creategrouppost.go` post creation
- `utils/extractphones.go` content loader

## 📞 Support & Customization

- **Email**: [haronkibetrutoh@gmail.com](mailto:haronkibetrutoh@gmail.com)
- **WhatsApp**: [+254718448461](https://wa.me/254718448461)

## ⚠️ Disclaimer

This tool is for educational and personal use only. Ensure compliance with Facebook’s Terms, applicable laws, and group rules.

---

Made with ❤️ using Go and Facebook GraphQL APIs
