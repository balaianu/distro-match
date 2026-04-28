// Wizard state management
export const wizardState = {
  currentStep: 0,
  preferences: {
    experienceLevel: null,
    useCase: null,
    hardware: {
      ram: null,
      disk: null,
      type: null,
    },
    desktopEnvironment: null,
    releaseModel: null,
    packageManager: null,
    supportLevel: null,
    philosophy: null,
    privacyLevel: null,
    learningGoal: null,
  },
};

export const questions = [
  {
    id: 'experienceLevel',
    question: 'What is your experience level with Linux?',
    options: [
      { value: 'beginner', label: 'Beginner - Never used Linux or very limited experience', icon: '🌱' },
      { value: 'intermediate', label: 'Intermediate - Some Linux experience, comfortable with basic commands', icon: '🌿' },
      { value: 'advanced', label: 'Advanced - Experienced user, comfortable with system administration', icon: '🌳' },
      { value: 'expert', label: 'Expert - Deep Linux knowledge, comfortable with complex configurations', icon: '🌲' },
    ],
    required: true,
  },
  {
    id: 'useCase',
    question: 'What will you primarily use this system for? (Select all that apply)',
    options: [
      { value: 'general_desktop', label: 'General Desktop - Web browsing, email, documents, media', icon: '🖥️' },
      { value: 'development', label: 'Development - Programming, coding, web development', icon: '💻' },
      { value: 'server', label: 'Server - Web servers, databases, containers, cloud', icon: '🖧' },
      { value: 'security', label: 'Security - Security research, penetration testing', icon: '🔒' },
      { value: 'gaming', label: 'Gaming - PC gaming with Steam/Proton', icon: '🎮' },
      { value: 'content_creation', label: 'Content Creation - Video editing, graphic design, audio', icon: '🎨' },
      { value: 'not_sure', label: 'Not Sure - Multiple use cases', icon: '🤔' },
    ],
    required: true,
    multiSelect: true,
  },
  {
    id: 'ram',
    question: 'How much RAM does your system have?',
    options: [
      { value: 'lt_2gb', label: 'Less than 2GB', icon: '📉' },
      { value: '2_4gb', label: '2GB - 4GB', icon: '📊' },
      { value: '4_8gb', label: '4GB - 8GB', icon: '📈' },
      { value: '8_16gb', label: '8GB - 16GB', icon: '🚀' },
      { value: 'gt_16gb', label: '16GB or more', icon: '⚡' },
      { value: 'not_sure', label: 'Not Sure', icon: '❓' },
    ],
    required: true,
    subStep: true,
    parent: 'hardware',
  },
  {
    id: 'disk',
    question: 'How much disk space is available?',
    options: [
      { value: 'lt_20gb', label: 'Less than 20GB', icon: '💾' },
      { value: '20_50gb', label: '20GB - 50GB', icon: '💿' },
      { value: '50_100gb', label: '50GB - 100GB', icon: '📀' },
      { value: 'gt_100gb', label: '100GB or more', icon: '💽' },
      { value: 'not_sure', label: 'Not Sure', icon: '❓' },
    ],
    required: true,
    subStep: true,
    parent: 'hardware',
  },
  {
    id: 'hardwareType',
    question: 'What type of hardware are you using?',
    options: [
      { value: 'desktop', label: 'Desktop PC - Tower computer', icon: '🖥️' },
      { value: 'laptop', label: 'Laptop - Portable computer', icon: '💻' },
      { value: 'server', label: 'Server - Headless machine', icon: '🖧' },
      { value: 'raspberry_pi', label: 'Raspberry Pi / ARM device', icon: '🍓' },
      { value: 'virtual_machine', label: 'Virtual Machine / Cloud', icon: '☁️' },
      { value: 'not_sure', label: 'Not Sure', icon: '❓' },
    ],
    required: true,
    subStep: true,
    parent: 'hardware',
  },
  {
    id: 'desktopEnvironment',
    question: 'Which desktop environment do you prefer? (Select all that apply)',
    options: [
      { value: 'no_preference', label: 'No Preference - Recommend based on other criteria', icon: '🎯' },
      { value: 'gnome', label: 'GNOME - Modern, workflow-focused', icon: '🔷' },
      { value: 'kde', label: 'KDE Plasma - Customizable, feature-rich', icon: '💠' },
      { value: 'xfce', label: 'XFCE - Lightweight, traditional', icon: '🔸' },
      { value: 'mate', label: 'MATE - Classic GNOME 2 style', icon: '🟢' },
      { value: 'cinnamon', label: 'Cinnamon - Traditional desktop', icon: '🟤' },
      { value: 'pantheon', label: 'Pantheon - macOS-like experience', icon: '🍎' },
      { value: 'i3', label: 'i3/Sway - Tiling window manager', icon: '🪟' },
      { value: 'lxqt', label: 'LXQt/LXDE - Very lightweight', icon: '🟣' },
    ],
    required: false,
    multiSelect: true,
  },
  {
    id: 'releaseModel',
    question: 'How do you prefer software updates?',
    options: [
      { value: 'no_preference', label: 'No Preference - Recommend based on experience level', icon: '🎯' },
      { value: 'stable_lts', label: 'Stable/LTS - Infrequent updates, maximum stability', icon: '🏛️' },
      { value: 'semi_rolling', label: 'Semi-Rolling - Regular updates with testing period', icon: '🔄' },
      { value: 'rolling', label: 'Rolling - Always latest software, cutting-edge', icon: '⚡' },
      { value: 'fixed_release', label: 'Fixed Release - Predictable 6-12 month schedule', icon: '📅' },
    ],
    required: false,
  },
  {
    id: 'packageManager',
    question: 'Which package manager do you prefer?',
    options: [
      { value: 'no_preference', label: 'No Preference - Recommend based on other criteria', icon: '🎯' },
      { value: 'apt', label: 'APT - Debian/Ubuntu family', icon: '📦' },
      { value: 'dnf', label: 'DNF/RPM - Fedora/RHEL family', icon: '📥' },
      { value: 'pacman', label: 'Pacman - Arch family', icon: '🎁' },
      { value: 'portage', label: 'Portage - Gentoo', icon: '⚙️' },
    ],
    required: false,
  },
  {
    id: 'supportLevel',
    question: 'What level of support do you need? (Select all that apply)',
    options: [
      { value: 'no_preference', label: 'No Preference - Any level acceptable', icon: '🎯' },
      { value: 'extensive', label: 'Extensive Community - Large forums, beginner-friendly', icon: '👥' },
      { value: 'professional', label: 'Professional Support - Paid enterprise support available', icon: '💼' },
      { value: 'documentation', label: 'Documentation - Comprehensive official docs', icon: '📚' },
      { value: 'minimal', label: 'Minimal - Comfortable figuring things out', icon: '🚶' },
    ],
    required: false,
    multiSelect: true,
  },
  {
    id: 'philosophy',
    question: 'Are there any philosophical requirements?',
    options: [
      { value: 'no_preference', label: 'No Preference - Pragmatic approach is fine', icon: '🎯' },
      { value: 'free_software', label: '100% Free Software - FSF-approved, no proprietary code', icon: '🆓' },
      { value: 'freedom', label: 'User Freedom - Flexible, customizable', icon: '🦅' },
      { value: 'corporate', label: 'Corporate-Backed - Prefer commercially supported distros', icon: '🏢' },
    ],
    required: false,
  },
  {
    id: 'privacyLevel',
    question: 'What level of privacy/security do you need?',
    options: [
      { value: 'casual', label: 'Casual - Standard privacy, no special requirements', icon: '🎯' },
      { value: 'enhanced', label: 'Enhanced - Private browsing, anti-tracking, encryption', icon: '🔐' },
      { value: 'high', label: 'High - Anonymity-focused, Tor/VPN, anti-forensics', icon: '🛡️' },
      { value: 'extreme', label: 'Extreme - Maximum isolation, compartmentalization, Qubes-level', icon: '🔒' },
      { value: 'not_sure', label: 'Not Sure', icon: '❓' },
    ],
    required: false,
  },
  {
    id: 'learningGoal',
    question: 'What is your primary goal with this Linux installation?',
    options: [
      { value: 'productivity', label: 'Productivity - Just get work done, minimal maintenance', icon: '⚡' },
      { value: 'balance', label: 'Balance - Learn while being productive', icon: '⚖️' },
      { value: 'learning', label: 'Learning - Want to understand Linux deeply', icon: '📚' },
      { value: 'not_sure', label: 'Not Sure', icon: '❓' },
    ],
    required: false,
  },
];

export function setPreference(key, value) {
  if (key === 'ram' || key === 'disk' || key === 'hardwareType') {
    if (key === 'hardwareType') {
      wizardState.preferences.hardware.type = value;
    } else {
      wizardState.preferences.hardware[key] = value;
    }
  } else {
    wizardState.preferences[key] = value;
  }
}

export function getPreference(key) {
  if (key === 'ram' || key === 'disk') {
    return wizardState.preferences.hardware[key];
  }
  if (key === 'hardwareType') {
    return wizardState.preferences.hardware.type;
  }
  return wizardState.preferences[key];
}

export function nextStep() {
  wizardState.currentStep++;
}

export function prevStep() {
  if (wizardState.currentStep > 0) {
    wizardState.currentStep--;
  }
}

export function resetWizard() {
  wizardState.currentStep = 0;
  wizardState.preferences = {
    experienceLevel: null,
    useCase: null,
    hardware: {
      ram: null,
      disk: null,
      type: null,
    },
    desktopEnvironment: null,
    releaseModel: null,
    packageManager: null,
    supportLevel: null,
    philosophy: null,
    privacyLevel: null,
    learningGoal: null,
  };
}

export function getCurrentQuestion() {
  return questions[wizardState.currentStep];
}

export function canSkipStep() {
  const currentQ = getCurrentQuestion();
  return !currentQ.required;
}

export function isStepComplete(stepIndex) {
  const q = questions[stepIndex];
  if (q.subStep) {
    if (q.id === 'hardwareType') {
      return wizardState.preferences.hardware.type !== null;
    }
    return wizardState.preferences.hardware[q.id] !== null;
  }
  return wizardState.preferences[q.id] !== null || !q.required;
}

export function getProgress() {
  const totalSteps = questions.length;
  const completedSteps = questions.filter((_, i) => isStepComplete(i)).length;
  return {
    current: wizardState.currentStep + 1,
    total: totalSteps,
    completed: completedSteps,
  };
}
