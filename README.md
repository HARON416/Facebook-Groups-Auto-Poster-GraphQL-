# 📱 Facebook Groups GraphQL Auto Poster

A powerful Go application that automatically posts content to Facebook groups using GraphQL APIs. This tool helps you manage and schedule posts across multiple Facebook groups with images and descriptions.

## ✨ Features

- **🔄 Automatic Group Discovery**: Fetches all your joined Facebook groups using GraphQL API
- **📸 Image Upload**: Automatically uploads images to Facebook and attaches them to posts
- **📝 Smart Content Management**: Reads item descriptions and images from organized directories
- **⏱️ Intelligent Scheduling**: Random delays between posts to avoid detection
- **🎯 Randomized Posting**: Shuffles groups and content for natural posting patterns
- **📊 Batch Processing**: Processes multiple groups and content items efficiently

## 🏗️ Architecture

The application is built with a modular architecture:

```
POST_WITH_GRAPHQL/
├── main.go                 # Main application entry point
├── utils/                  # Utility modules
│   ├── fetchgroups.go      # Facebook groups fetching via GraphQL
│   ├── creategrouppost.go  # Post creation and submission
│   ├── uploadimage.go      # Image upload functionality
│   ├── extractphones.go    # Content extraction from directories
│   ├── returnrandomnumber.go # Random delay generation
│   └── buildattachments.go # Attachment building utilities
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
   itemsPath := "/path/to/your/content/directory"
   ```

2. **Ensure your Facebook session is active** (the application uses Facebook's GraphQL API)

## 📖 Usage

### Running the Application

```bash
go run main.go
```

### How It Works

1. **Group Discovery**: The app fetches all your joined Facebook groups
2. **Content Loading**: Reads item descriptions and images from your content directory
3. **Randomization**: Shuffles groups and content for natural posting
4. **Image Upload**: Uploads images to Facebook and gets photo IDs
5. **Post Creation**: Creates posts with descriptions and attached images
6. **Scheduling**: Adds random delays between posts (2-5 minutes)

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

## 🔧 Technical Details

### GraphQL Integration

The application uses Facebook's GraphQL API for:

- **Group Fetching**: `GroupsCometAllJoinedGroupsSectionPaginationQuery`
- **Post Creation**: Custom mutation for group posts
- **Image Upload**: Facebook's media upload endpoints

### Key Components

- **`FetchGroups()`**: Retrieves all joined Facebook groups with pagination
- **`ExtractItems()`**: Reads and processes content from directory structure
- **`UploadImage()`**: Handles image upload to Facebook
- **`CreateGroupPost()`**: Creates and submits posts to groups
- **`ReturnRandomNumberBetween2And5()`**: Generates random delays

### Error Handling

The application includes comprehensive error handling for:

- Network connectivity issues
- Facebook API rate limiting
- File system errors
- Invalid content formats

## 🛡️ Safety Features

- **Random Delays**: 2-5 minute intervals between posts
- **Content Shuffling**: Randomizes both groups and content order
- **Error Recovery**: Continues processing even if individual posts fail
- **Session Management**: Uses Facebook's session tokens for authentication

## 📋 Requirements

- **Go Version**: 1.24.5+
- **Dependencies**:
  - `github.com/klauspost/compress v1.18.0`
- **Facebook Account**: Active session with joined groups
- **Content**: Organized directory with images and descriptions

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
