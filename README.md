# 📱 Facebook Groups GraphQL Auto Poster

A powerful Go application that automatically posts content to Facebook groups using GraphQL APIs. This tool helps you manage and schedule posts across multiple Facebook groups with images and descriptions.

## ✨ Features

- **🔄 Automatic Group Discovery**: Fetches all your joined Facebook groups using GraphQL API with pagination support
- **📸 Image Upload**: Automatically uploads images to Facebook and attaches them to posts
- **📝 Smart Content Management**: Reads item descriptions and images from organized directories
- **⏱️ Intelligent Scheduling**: Random delays between posts to avoid detection (2-4 minutes)
- **🎯 Randomized Posting**: Shuffles groups and content for natural posting patterns
- **📊 Batch Processing**: Processes multiple groups and content items efficiently
- **🔐 Dynamic Credential Management**: Automatic cookie and session token updates from curl requests

## 🏗️ Architecture

The application is built with a modular architecture:

```
POST_WITH_GRAPHQL/
├── main.go                 # Main application entry point
├── go.mod                  # Go module dependencies
├── go.sum                  # Dependency checksums
├── makefile                # Build and run commands
├── cookies/                # Cookie storage directory
│   ├── fetchgroupscookies.txt
│   ├── uploadphotocookies.txt
│   └── creategrouppostcookies.txt
└── utils/                  # Utility modules
    ├── fetchgroups.go      # Facebook groups fetching via GraphQL
    ├── creategrouppost.go  # Post creation and submission
    ├── uploadimage.go      # Image upload functionality
    ├── extractphones.go    # Content extraction from directories
    ├── returnrandomnumber.go # Random delay generation
    ├── buildattachments.go # Attachment building utilities
    ├── updatefetchgroupscookies.go # Cookie update for group fetching
    ├── updateuploadphotocookies.go # Cookie update for image upload
    └── updatecreategrouppostcookies.go # Cookie update for post creation
```

## 🚀 Getting Started

### Prerequisites

- Go 1.24.5 or higher
- Facebook account with joined groups
- Organized content directory structure

### Installation

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd POST_WITH_GRAPHQL
   ```

2. **Install dependencies**

   ```bash
   go mod tidy
   ```

3. **Set up your content directory**
   ```
   /path/to/your/content/
   ├── item1/
   │   ├── details.txt
   │   ├── image1.jpg
   │   └── image2.jpg
   ├── item2/
   │   ├── details.txt
   │   └── image1.jpg
   └── ...
   ```

### Configuration

1. **Update the content path** in `main.go`:

   ```go
   phonesPath := "/path/to/your/content/directory"
   ```

2. **Set up cookie files** in the `cookies/` directory:
   - `fetchgroupscookies.txt` - Curl request for fetching groups
   - `uploadphotocookies.txt` - Curl request for uploading images
   - `creategrouppostcookies.txt` - Curl request for creating posts

## 🔐 Cookie Management

### Cookie Files Structure

The application uses three cookie files located in the `cookies/` directory:

```
cookies/
├── fetchgroupscookies.txt      # For fetching Facebook groups
├── uploadphotocookies.txt      # For uploading images to Facebook
└── creategrouppostcookies.txt  # For creating posts in groups
```

### How to Update Cookie Files

#### **Step 1: Get Fresh Curl Requests**

1. **Open Facebook in your browser** and log in
2. **Open Developer Tools** (F12) and go to the Network tab
3. **Perform the action** you want to capture:
   - **For Groups**: Go to your Facebook groups page
   - **For Image Upload**: Try to upload an image to a group
   - **For Post Creation**: Try to create a post in a group
4. **Find the relevant request** in the Network tab
5. **Right-click the request** → Copy → Copy as cURL

#### **Step 2: Update the Cookie Files**

**For Fetching Groups (`fetchgroupscookies.txt`):**

1. **Login to Facebook**
2. **Go to "My Groups"**
3. **Sort alphabetically**
4. **Copy curl request in Network tab**

**For Image Upload (`uploadphotocookies.txt`):**

1. **Go to any of your groups**
2. **Trigger post creation dialog**
3. **Select images**
4. **Copy curl request in Network tab**

**For Post Creation (`creategrouppostcookies.txt`):**

1. **Go to any of your groups**
2. **Trigger post creation dialog**
3. **Input your post and submit**
4. **Copy curl request in Network tab**

#### **Step 3: Run the Update Functions**

After updating the cookie files, run the corresponding update functions:

```go
// Update group fetching credentials
err := utils.UpdateFetchGroupsCookies()

// Update image upload credentials
err := utils.UpdateUploadPhotoCookies()

// Update post creation credentials
err := utils.UpdateCreateGroupPostCookies()
```

### Cookie File Format

Each cookie file should contain a complete cURL request. Example format:

```bash
curl 'https://www.facebook.com/api/graphql/' \
  -H 'authority: www.facebook.com' \
  -H 'accept: */*' \
  -H 'accept-language: en-US,en;q=0.9' \
  -H 'content-type: application/x-www-form-urlencoded' \
  -H 'cookie: your_cookies_here' \
  -H 'origin: https://www.facebook.com' \
  -H 'referer: https://www.facebook.com/groups/' \
  -H 'user-agent: your_user_agent' \
  --data-raw 'your_form_data_here'
```

### Automatic Cookie Updates

The application automatically updates credentials when you run it:

```go
// These functions are called automatically in main.go
err := utils.UpdateFetchGroupsCookies()
err := utils.UpdateUploadPhotoCookies()
err := utils.UpdateCreateGroupPostCookies()
```

### Verification

The update functions include verification to ensure cookies were updated correctly:

- ✅ Checks that critical parameters were updated
- ✅ Verifies headers were updated properly
- ✅ Confirms cookies were set correctly
- ✅ Shows debug output of what was updated

### Troubleshooting Cookie Updates

**If cookie updates fail:**

1. **Check file format**: Ensure the cURL request is complete and properly formatted
2. **Verify file location**: Make sure files are in the `cookies/` directory
3. **Check permissions**: Ensure the application can read the cookie files
4. **Update frequency**: Facebook cookies expire regularly, update them every few hours

**Common issues:**

- ❌ **Empty cookie files**: Make sure files contain valid cURL requests
- ❌ **Expired cookies**: Facebook cookies expire quickly, update frequently
- ❌ **Wrong request type**: Ensure you're copying the correct GraphQL request

### Security Notes

- 🔒 **Keep cookie files private**: Don't commit them to version control
- 🔒 **Regular updates**: Facebook cookies expire frequently
- 🔒 **Session management**: Use fresh browser sessions for new cookies
- 🔒 **Account security**: Don't share cookie files with others

## 📖 Usage

### Running the Application

```bash
# Using make
make run

# Or directly with go
go run main.go
```

### How It Works

1. **Credential Updates**: Automatically updates cookies and session tokens from curl requests
2. **Group Discovery**: Fetches all your joined Facebook groups with pagination support
3. **Content Loading**: Reads item descriptions and images from your content directory
4. **Randomization**: Shuffles groups and content for natural posting patterns
5. **Image Upload**: Uploads images to Facebook and gets photo IDs
6. **Post Creation**: Creates posts with descriptions and attached images
7. **Scheduling**: Adds random delays between posts (2-4 minutes)

### Content Format

Your content directory should follow this structure:

```
content/
├── item1/
│   ├── details.txt          # Contains: "description: Your post description here..."
│   ├── image1.jpg
│   └── image2.jpg
├── item2/
│   ├── details.txt
│   └── image1.jpg
└── ...
```

The `details.txt` file should contain:

```
description: Your post description with features...separated by ellipsis...for better formatting
```

**Note**: The application automatically:

- Converts descriptions to uppercase
- Adds ✅ emoji to each feature
- Splits content by "..." and formats with line breaks

## 🔧 Technical Details

### GraphQL Integration

The application uses Facebook's GraphQL API for:

- **Group Fetching**: `GroupsCometAllJoinedGroupsSectionPaginationQuery` with pagination
- **Post Creation**: Custom mutation for group posts with comprehensive input variables
- **Image Upload**: Facebook's media upload endpoints

### 🔄 Credential Management

The application includes a sophisticated dynamic credential update system:

#### **Automatic Credential Updates**

```go
// Update credentials from curl request files
err := utils.UpdateFetchGroupsCookies()
err := utils.UpdateUploadPhotoCookies()
err := utils.UpdateCreateGroupPostCookies()
```

#### **Workflow**:

1. **Add new curl request** to appropriate cookie file in `cookies/` directory
2. **Run** the corresponding update function
3. **All credentials updated** automatically!

#### **Features**:

- ✅ **Dynamic updates** - No hardcoded credentials
- ✅ **Comprehensive** - Updates all parameters, headers, and cookies
- ✅ **URL decoding** - Handles encoded values properly
- ✅ **Verification** - Checks that updates worked correctly
- ✅ **Debug output** - Shows what was updated

### Key Components

- **`FetchGroups()`**: Retrieves all joined Facebook groups with pagination support
- **`ExtractPhones()`**: Reads and processes content from directory structure
- **`UploadImage()`**: Handles image upload to Facebook
- **`CreateGroupPost()`**: Creates and submits posts to groups with comprehensive input variables
- **`ReturnRandomNumber()`**: Generates random delays between 2-4 minutes

### Error Handling

The application includes comprehensive error handling for:

- Network connectivity issues
- Facebook API rate limiting
- File system errors
- Invalid content formats
- Cookie update failures

## 🛡️ Safety Features

- **Random Delays**: 2-4 minute intervals between posts
- **Content Shuffling**: Randomizes both groups and content order
- **Error Recovery**: Continues processing even if individual posts fail
- **Session Management**: Uses Facebook's session tokens for authentication
- **Credential Rotation**: Automatic cookie and session token updates

## 📋 Requirements

- **Go Version**: 1.24.5+
- **Dependencies**:
  - `github.com/klauspost/compress v1.18.0`
- **Facebook Account**: Active session with joined groups
- **Content**: Organized directory with images and descriptions
- **Cookie Files**: Valid curl requests in `cookies/` directory

## 🛠️ Development

### Build Commands

```bash
# Run the application
make run

# Tidy dependencies
make tidy

# Build binary
go build -o facebook-poster main.go
```

### Project Structure

- **Main Logic**: `main.go` orchestrates the entire workflow
- **GraphQL Operations**: `utils/fetchgroups.go` handles group discovery
- **Content Processing**: `utils/extractphones.go` manages content extraction
- **Post Creation**: `utils/creategrouppost.go` handles post submission
- **Image Upload**: `utils/uploadimage.go` manages media uploads
- **Credential Management**: Cookie update utilities in `utils/`

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Support & Customization

Need help with setup, customization, or have questions? Get in touch:

- **📧 Email**: [haronkibetrutoh@gmail.com](mailto:haronkibetrutoh@gmail.com)
- **📱 WhatsApp**: [+254718448461](https://wa.me/254718448461)

I'm available for:

- Setup assistance and troubleshooting
- Custom feature development
- Integration with other platforms
- Bulk customization requests
- Training and consultation

## ⚠️ Disclaimer

This tool is for educational and personal use only. Please ensure you comply with:

- Facebook's Terms of Service
- Applicable laws and regulations
- Group-specific posting guidelines

Use responsibly and respect community guidelines when posting to Facebook groups.

---

**Made with ❤️ using Go and Facebook GraphQL APIs**
