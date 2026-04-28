# Contributing to DistroMatch

Thank you for your interest in contributing to DistroMatch! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on what is best for the community
- Show empathy towards other community members

## Getting Started

### Prerequisites

- Node.js 18.x or higher
- npm or yarn
- Git
- Basic knowledge of HTML, CSS, and JavaScript

### Initial Setup

1. Fork the repository
2. Clone your fork:
```bash
git clone https://github.com/balaianu/distro-match.git
cd distro-match
```

3. Install dependencies:
```bash
npm install
```

4. Start development server:
```bash
npm run dev
```

5. Open `http://localhost:4321` in your browser

## How to Contribute

### Areas for Contribution

1. **Add Linux Distributions** - Add new distributions to the database or update existing information
2. **Improve the Algorithm** - Enhance scoring logic, add new criteria, improve weighting
3. **Enhance UI/UX** - Improve visual design, accessibility, animations, mobile experience
4. **Fix Bugs** - Fix reported issues, improve error handling, edge cases
5. **Documentation** - Improve documentation, add examples, translate
6. **Internationalization** - Add translations, support RTL languages
7. **Testing** - Write unit tests, integration tests, improve coverage

## Development Workflow

1. **Create a Branch**
```bash
git checkout -b feature/your-feature-name
```

2. **Make Changes** - Edit source files, test locally, ensure build succeeds

3. **Test Your Changes**
```bash
npm run build
npm run preview
```

4. **Commit Changes**
```bash
git add .
git commit -m "Brief description of changes"
```

5. **Push to Your Fork**
```bash
git push origin feature/your-feature-name
```

6. **Create Pull Request** - Go to the original repository, click "New Pull Request", describe your changes

### Commit Message Guidelines

- Use present tense ("Add feature" not "Added feature")
- Use imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit to 72 characters for the first line
- Reference issues in the commit message body

## Adding Linux Distributions

### Location

Edit `src/data/distros.json`

### Schema

Each distribution must follow this schema:

```json
{
  "id": "unique-identifier",
  "name": "Distribution Name",
  "version": "Version Number",
  "experience_level": ["beginner", "intermediate", "advanced", "expert"],
  "use_cases": ["general_desktop", "development", "server", "security", "gaming", "content_creation", "old_hardware", "privacy"],
  "min_ram_gb": 2,
  "recommended_ram_gb": 4,
  "min_disk_gb": 10,
  "cpu_architecture": ["x86_64", "arm64", "i386"],
  "desktop_environments": ["GNOME", "KDE Plasma", "XFCE", "MATE", "Cinnamon", "Pantheon", "Any"],
  "release_model": "stable_lts",
  "package_manager": "apt",
  "community_support": "extensive",
  "professional_support": true,
  "documentation_quality": "high",
  "philosophy": "pragmatic",
  "based_on": "debian",
  "description": "Brief description of the distribution",
  "official_website": "https://example.com",
  "download_page": "https://example.com/download",
  "strengths": ["Strength 1", "Strength 2", "Strength 3"],
  "weaknesses": ["Weakness 1", "Weakness 2"],
  "similar_distros": ["similar-distro-1", "similar-distro-2"]
}
```

### Guidelines

- **ID**: Use lowercase with hyphens (e.g., "ubuntu-24-04")
- **Experience Level**: Must include at least one valid level
- **Use Cases**: Must include at least one valid use case
- **Hardware Requirements**: Be realistic and accurate
- **Desktop Environments**: List actual supported DEs
- **Release Model**: Use one of the valid values
- **Package Manager**: Use the actual package manager
- **Support**: Be accurate about community and professional support
- **Philosophy**: Use one of the valid philosophy values
- **Description**: Keep concise but informative (2-3 sentences)
- **Strengths**: List 3-5 key strengths
- **Weaknesses**: List 1-3 notable weaknesses
- **Similar Distros**: List 2-3 similar distributions

## Modifying the Wizard

### Adding Questions

Edit `src/scripts/wizard-state.js`

Add to the `questions` array:

```javascript
{
  id: 'newQuestionId',
  question: 'What is your question?',
  options: [
    { value: 'option1', label: 'Option 1 description' },
    { value: 'option2', label: 'Option 2 description' },
  ],
  required: true,
  subStep: false,
  parent: null,
}
```

### Guidelines for Questions

- **ID**: Use camelCase, descriptive
- **Question**: Clear and concise
- **Options**: Provide 3-6 options when possible
- **Required**: Mark as required if critical for recommendations
- **Sub-step**: Use for related questions (e.g., hardware specs)
- **Labels**: Be descriptive and helpful

## Improving the Algorithm

### Location

The scoring algorithm is in `src/scripts/recommendation-algorithm.js`

### Adding New Scoring Criteria

1. Add scoring function
2. Add to WEIGHTS object
3. Add to calculateScores function

### Guidelines

- **Score Range**: Always return 0-1
- **Weights**: Ensure all weights sum to 1.0
- **Default Handling**: Return 1 for null/no_preference
- **Consistency**: Keep scoring logic consistent across similar criteria

## Style Guide

### Code Style

- **JavaScript**: Use modern ES6+ syntax
- **Indentation**: 2 spaces
- **Quotes**: Use single quotes for strings
- **Semicolons**: Use semicolons
- **Naming**: camelCase for variables/functions, PascalCase for components

### Astro/HTML Style

- Use semantic HTML elements
- Follow Astro best practices
- Use Tailwind classes for styling
- Keep components focused and small

### CSS Style

- Prefer Tailwind utility classes
- Use custom classes only for complex components
- Follow mobile-first responsive design
- Use semantic class names

## Testing

### Manual Testing Checklist

Before submitting a PR:

- [ ] Build succeeds: `npm run build`
- [ ] Preview works: `npm run preview`
- [ ] Wizard flow completes without errors
- [ ] All questions display correctly
- [ ] Results display correctly
- [ ] Responsive design works on mobile
- [ ] Keyboard navigation works
- [ ] No console errors

## Submitting Changes

### Pull Request Checklist

- [ ] Code follows style guide
- [ ] Build succeeds without errors
- [ ] Manual testing completed
- [ ] Documentation updated (if needed)
- [ ] Commit messages are clear
- [ ] PR description explains changes
- [ ] Related issues are referenced

### PR Description Template

```markdown
## Description
Brief description of what this PR does.

## Changes
- List of main changes
- Files modified
- New features added

## Testing
- How you tested the changes
- Test results

## Related Issues
Closes #123
```

## Reporting Issues

### Bug Reports

When reporting a bug, include:

- **Description**: What happened
- **Steps to Reproduce**: How to reproduce the issue
- **Expected Behavior**: What should happen
- **Actual Behavior**: What actually happened
- **Environment**: Browser, OS, Node version
- **Screenshots**: If applicable

### Feature Requests

When requesting a feature, include:

- **Description**: What you want
- **Use Case**: Why you need it
- **Alternatives**: What you've considered
- **Additional Context**: Any other relevant information

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to DistroMatch!
