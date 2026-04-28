# DistroMatch

<div align="center">

**Find your perfect Linux distribution**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Astro](https://img.shields.io/badge/Astro-FF5D01?style=flat&logo=astro&logoColor=white)](https://astro.build)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind-06B6D4?style=flat&logo=tailwindcss&logoColor=white)](https://tailwindcss.com)

[Live Demo](https://distro-match.com) · [Report Bug](https://github.com/balaianu/distro-match/issues) · [Request Feature](https://github.com/balaianu/distro-match/issues)

</div>

A guided wizard that helps users find the perfect Linux distribution based on their experience level, use case, hardware specifications, and personal preferences. Built with Astro and Tailwind CSS for a modern, responsive, and fast experience.

## ✨ Features

- **Smart Recommendations** - Weighted scoring algorithm matches users to suitable Linux distributions
- **Guided Wizard** - Step-by-step questionnaire with 10 decision points
- **Hardware-Aware** - Filters distros based on RAM, CPU architecture, and disk space
- **Distro Database** - Comprehensive database of popular Linux distributions
- **Responsive Design** - Works seamlessly on desktop, tablet, and mobile devices
- **Static Site** - Zero server costs, instant deployment, CDN delivery
- **Open Source** - Fully MIT licensed

## 🚀 Quick Start

### Prerequisites

- Node.js 18.x or higher
- npm or yarn

### Installation

```bash
# Clone the repository
git clone https://github.com/balaianu/distro-match.git
cd distro-match

# Install dependencies
npm install

# Start development server
npm run dev
```

Open [http://localhost:4321](http://localhost:4321) in your browser.

### Building for Production

```bash
npm run build
```

The built files will be in the `dist/` directory.

## � Usage

1. **Start the Wizard** - Click "Find Your Perfect Match" to begin
2. **Answer Questions** - Provide your preferences across 10 categories
3. **Get Recommendations** - Receive personalized distro recommendations with match percentages
4. **Explore Details** - View detailed information about each recommended distribution

## 🏗️ Project Structure

```
distro-match/
├── src/
│   ├── data/
│   │   └── distros.json          # Linux distribution database
│   ├── layouts/
│   │   └── Layout.astro          # Main layout component
│   ├── pages/
│   │   ├── index.astro           # Main wizard page
│   │   ├── distros.astro         # Distros exploration page
│   │   └── about.astro           # About page
│   ├── scripts/
│   │   ├── wizard-state.js       # Wizard state management
│   │   └── recommendation-algorithm.js  # Scoring algorithm
│   └── styles/
│       └── global.css            # Global styles
├── docs/                         # Documentation
├── public/                       # Static assets
├── astro.config.mjs             # Astro configuration
├── tailwind.config.mjs          # Tailwind configuration
└── package.json                 # Dependencies
```

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guidelines](docs/CONTRIBUTING.md) for details on how to contribute.

### Quick Contribution Steps

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 Adding Linux Distributions

To add a new Linux distribution to the database:

1. Edit `src/data/distros.json`
2. Add a new object following the existing schema
3. Rebuild: `npm run build`

See [Contributing Guidelines](docs/CONTRIBUTING.md) for detailed schema information.

## 🌐 Deployment

### Cloudflare Pages (Recommended)

1. Connect your GitHub repository to Cloudflare Pages
2. Build command: `npm run build`
3. Output directory: `dist`
4. Add custom domain: `distro-match.com`

### Other Options

- **Netlify** - Build command: `npm run build`, Publish directory: `dist`
- **Vercel** - Build command: `npm run build`, Output directory: `dist`
- **GitHub Pages** - Build and deploy `dist/` directory

## 🛠️ Tech Stack

- **Framework**: [Astro](https://astro.build) 5.x
- **Styling**: [Tailwind CSS](https://tailwindcss.com) 3.x
- **JavaScript**: Vanilla ES6+
- **Architecture**: Static HTML/CSS/JS with client-side processing

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Linux distributions data sourced from official project websites
- Built with [Astro](https://astro.build) and [Tailwind CSS](https://tailwindcss.com)
- Design inspiration from DistroWatch and various Linux comparison sites

## 💬 Support

- **Issues**: [Report bugs or request features](https://github.com/balaianu/distro-match/issues)
- **Discussions**: [Ask questions or share ideas](https://github.com/balaianu/distro-match/discussions)

## ☕ Support the Project

If you find this project helpful, consider supporting its development:

[![Buy Me A Coffee](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/balaianu)

---

<div align="center">

Made with ❤️ by [balaianu](https://github.com/balaianu)

</div>
