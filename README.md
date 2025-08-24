# Facebook Groups Auto Poster

A powerful Go application that automatically posts content with images to multiple Facebook groups. This tool extracts group information, uploads images, and creates posts across all your joined Facebook groups with intelligent randomization and progress tracking.

## 🚀 Features

- **Automatic Group Discovery**: Extracts all Facebook groups you've joined
- **Image Upload**: Supports multiple image uploads per post
- **Smart Content Management**: Randomly selects and shuffles content from local directories
- **Progress Tracking**: Real-time progress monitoring with success rate statistics
- **Rate Limiting**: Built-in delays to respect Facebook's rate limits
- **Dynamic Configuration**: Uses curl commands for easy authentication token updates
- **Comprehensive Logging**: Detailed logs for debugging and monitoring

## 📋 Prerequisites

- Go 1.21 or higher
- Active Facebook account with joined groups
- Local directory with content items (images and descriptions)

## 🛠️ Installation

1. **Clone the repository:**

```bash
git clone <repository-url>
cd facebook-groups-autoposter
```

2. **Install dependencies:**

```bash
go mod tidy
```

3. **Build the application:**

```bash
make build
# or manually: go build -o autoposter main.go
```

## 📁 Project Structure

```
AUTOPOSTER/
├── main.go                 # Main application entry point
├── utils/
│   ├── fetchgroups.go      # Facebook groups extraction
│   ├── uploadimage.go      # Image upload functionality
│   ├── createpost.go       # Post creation logic
│   └── extractitems.go     # Local content extraction
├── go.mod                  # Go module dependencies
├── makefile               # Build automation
└── README.md              # This file
```

## 🔧 Configuration

### 1. Content Directory Setup

Create a directory structure for your content:

```
/path/to/your/content/
├── item1/
│   ├── details.txt        # Post description/text
│   ├── image1.jpg
│   ├── image2.png
│   └── image3.jpeg
├── item2/
│   ├── details.txt
│   └── photo.jpg
└── ...
```

**Example `details.txt`:**

```
Get your dream SAMSUNG GALAXY NOTE10 (Ex UK) today with our Lipa Mdogo Mdogo plan!
Only KSh 9,199 deposit + KSh 1,240 weekly for 52 weeks.
12GB RAM/256GB storage.
Visit us at Pioneer Building, Kimathi Street, or call/WhatsApp 0718448461.
```

### 2. Update Content Path

Edit `main.go` line 35 to point to your content directory:

```go
itemsPath := "/path/to/your/content"
```

### 3. Authentication Setup

The application uses curl commands for authentication. You'll need to update three curl commands:

#### A. Groups Fetching (in `utils/fetchgroups.go`)

1. Open Facebook Groups page in your browser
2. Open Developer Tools (F12) → Network tab
3. Refresh the page and find the GraphQL request for groups
4. Right-click → Copy as cURL
5. Replace the `fetchGroupsCurl` variable in `utils/fetchgroups.go`

#### B. Image Upload (in `utils/uploadimage.go`)

1. Try uploading an image in Facebook
2. Capture the upload request from Network tab
3. Replace the `uploadImageCurl` variable in `utils/uploadimage.go`

#### C. Post Creation (in `utils/createpost.go`)

1. Create a test post in a Facebook group
2. Capture the post creation request
3. Replace the `createPostCurl` variable in `utils/createpost.go`

## 🚀 Usage

### Basic Usage

```bash
# Run the application
./autoposter

# Or with Go directly
go run main.go
```

### Command Line Options

Currently, the application runs with default settings. Configuration is done through code modification.

### Sample Output

```
🔄 Updating configuration from curl command...
✅ Configuration updated successfully!

📊 Groups extraction completed!
   Total groups found: 80
   Groups with valid data: 80

💾 JSON data ready (length: 15420 bytes)
🎉 All done!
📊 Total groups extracted: 80

📋 Processing group 1/80: Tech Enthusiasts Kenya (653028063286544)

🖼️ Starting image upload for: /path/to/content/item1/image1.jpg
✅ Image upload completed successfully!
   📷 Photo ID: 122229777992128715
   🔗 Image URL: https://scontent.xx.fbcdn.net/...

🎉 Post created successfully!
   📝 Post ID: unknown
   👥 Group ID: 653028063286544
   📄 Text: Get your dream SAMSUNG GALAXY NOTE10...
   🖼️  Photos: 3 images
   📊 Progress: 1/80 posts completed (1.3%)

[Waiting 2 minutes before next post...]

📋 Processing group 2/80: Mobile Phone Dealers...

🎯 Final Summary:
   ✅ Successful posts: 78/80
   📈 Success rate: 97.5%
   🏁 Auto-posting completed!
```

## ⚙️ Advanced Configuration

### Rate Limiting

The application includes a 2-minute delay between posts by default. Modify this in `main.go`:

```go
time.Sleep(2 * time.Minute) // Change duration as needed
```

### Image Selection

Images are randomly shuffled for each group. The application uses all images found in each content item directory.

### Error Handling

The application continues processing even if individual posts fail, providing detailed error logging for troubleshooting.

## 🔍 Troubleshooting

### Common Issues

1. **Authentication Errors**

   - **Problem**: `error:1357001` or "Login to continue"
   - **Solution**: Update curl commands with fresh authentication tokens

2. **No Groups Found**

   - **Problem**: Groups extraction returns 0 groups
   - **Solution**: Check Facebook Groups curl command and ensure you're logged in

3. **Image Upload Fails**

   - **Problem**: Upload returns authentication errors
   - **Solution**: Update image upload curl command with current session

4. **Post Creation Fails**
   - **Problem**: Posts return server errors
   - **Solution**: Update post creation curl command

### Debug Mode

For detailed debugging, the application provides comprehensive logging. Check the console output for:

- Request debugging information
- Response headers and bodies
- Form data parameters
- Authentication token status

### Getting Fresh Curl Commands

1. Open Facebook in your browser
2. Open Developer Tools (F12)
3. Go to Network tab
4. Perform the action (view groups, upload image, create post)
5. Find the corresponding request
6. Right-click → Copy as cURL (bash)
7. Replace the respective curl variable in the code

## 📊 Monitoring

The application provides real-time monitoring:

- **Progress Tracking**: Shows current group being processed
- **Success Rate**: Displays percentage of successful posts
- **Error Logging**: Detailed error messages for failed operations
- **Final Summary**: Complete statistics at the end

## ⚠️ Important Notes

### Rate Limiting

- Facebook has strict rate limits for posting
- The application includes 2-minute delays between posts
- Respect Facebook's terms of service

### Authentication

- Curl commands contain session tokens that expire
- Update authentication regularly (daily or as needed)
- Never share your curl commands (they contain private tokens)

### Content Guidelines

- Ensure your content complies with Facebook's community standards
- Avoid spam-like behavior
- Post relevant content to appropriate groups

## 📄 License

This project is for educational purposes. Please ensure compliance with Facebook's terms of service and applicable laws.

## 🆘 Support

For issues and questions:

1. Check the troubleshooting section above
2. Review the console output for error details
3. Ensure authentication tokens are current
4. Verify content directory structure

### Professional Setup Support

Need help setting up or configuring the auto-poster? I offer personalized setup assistance:

📱 **WhatsApp**: [+254718448461](https://wa.me/254718448461)  
📧 **Email**: [haronkibetrutoh@gmail.com](mailto:haronkibetrutoh@gmail.com)

**Setup Services Available:**

- Complete installation and configuration
- Authentication setup and troubleshooting
- Content directory optimization
- Custom modifications and enhancements
- Ongoing maintenance and updates

_Professional setup support available at a reasonable fee. Contact me for a quote!_

---

**⚠️ Disclaimer**: Use this tool responsibly and in accordance with Facebook's terms of service. The authors are not responsible for any account restrictions or violations resulting from misuse.
