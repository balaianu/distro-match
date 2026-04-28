import distros from '../data/distros.json';

// Scoring weights
const WEIGHTS = {
  experienceLevel: 0.25,
  useCase: 0.20,
  hardware: 0.20,
  hardwareType: 0.05,
  desktopEnvironment: 0.08,
  releaseModel: 0.03,
  packageManager: 0.03,
  privacyLevel: 0.03,
  learningGoal: 0.02,
  supportLevel: 0.02,
  philosophy: 0.01,
};

// Helper functions
function scoreExperienceLevel(distro, preference) {
  if (!preference || preference === 'not_sure') return 1;
  
  const experienceMap = {
    beginner: ['beginner'],
    intermediate: ['beginner', 'intermediate'],
    advanced: ['intermediate', 'advanced'],
    expert: ['advanced', 'expert'],
  };
  
  const acceptableLevels = experienceMap[preference] || [];
  return distro.experience_level.some(level => acceptableLevels.includes(level)) ? 1 : 0.3;
}

function scoreUseCase(distro, preference) {
  if (!preference || preference === 'not_sure') return 1;

  const useCaseMap = {
    general_desktop: ['general_desktop'],
    development: ['general_desktop', 'development'],
    server: ['server'],
    security: ['security'],
    gaming: ['general_desktop', 'gaming'],
    content_creation: ['general_desktop', 'content_creation'],
  };

  const acceptableUseCases = useCaseMap[preference] || [preference];
  return distro.use_cases.some(uc => acceptableUseCases.includes(uc)) ? 1 : 0.2;
}

function scoreHardware(distro, hardware) {
  let score = 1;

  if (hardware.ram && hardware.ram !== 'not_sure') {
    const ramMap = {
      lt_2gb: 0,
      '2_4gb': 2,
      '4_8gb': 4,
      '8_16gb': 8,
      gt_16gb: 16,
    };
    const requiredRam = ramMap[hardware.ram] || 4;
    if (distro.min_ram_gb > requiredRam) {
      score *= 0.2;
    }
  }

  // Infer CPU architecture from hardware type
  if (hardware.type && hardware.type !== 'not_sure') {
    const typeToArch = {
      desktop: ['x86_64'],
      laptop: ['x86_64'],
      server: ['x86_64'],
      raspberry_pi: ['arm64', 'armhf'],
      virtual_machine: ['x86_64'],
    };
    const acceptableArch = typeToArch[hardware.type] || ['x86_64'];
    if (!distro.cpu_architecture.some(arch => acceptableArch.includes(arch))) {
      score *= 0.1;
    }
  }

  if (hardware.disk && hardware.disk !== 'not_sure') {
    const diskMap = {
      lt_20gb: 0,
      '20_50gb': 20,
      '50_100gb': 50,
      gt_100gb: 100,
    };
    const requiredDisk = diskMap[hardware.disk] || 20;
    if (distro.min_disk_gb > requiredDisk) {
      score *= 0.2;
    }
  }

  return score;
}

function scoreDesktopEnvironment(distro, preference) {
  if (!preference || preference === 'no_preference') return 1;
  
  const deMap = {
    gnome: ['GNOME'],
    kde: ['KDE Plasma'],
    xfce: ['XFCE'],
    mate: ['MATE'],
    cinnamon: ['Cinnamon'],
    pantheon: ['Pantheon'],
    i3: ['i3', 'Sway'],
    lxqt: ['LXQt', 'LXDE'],
  };
  
  const acceptableDEs = deMap[preference] || [];
  if (acceptableDEs.length === 0) return 1;
  
  return distro.desktop_environments.some(de => acceptableDEs.includes(de)) ? 1 : 0.4;
}

function scoreReleaseModel(distro, preference) {
  if (!preference || preference === 'no_preference') return 1;
  
  const modelMap = {
    stable_lts: ['stable_lts'],
    semi_rolling: ['semi_rolling'],
    rolling: ['rolling', 'semi_rolling'],
    fixed_release: ['fixed_release', 'stable_lts'],
  };
  
  const acceptableModels = modelMap[preference] || [];
  return acceptableModels.includes(distro.release_model) ? 1 : 0.5;
}

function scorePackageManager(distro, preference) {
  if (!preference || preference === 'no_preference') return 1;

  const pmMap = {
    apt: ['apt'],
    dnf: ['dnf'],
    pacman: ['pacman'],
    portage: ['portage'],
  };

  const acceptablePMs = pmMap[preference] || [];
  return acceptablePMs.includes(distro.package_manager) ? 1 : 0.5;
}

function scoreSupportLevel(distro, preference) {
  if (!preference || preference === 'no_preference') return 1;
  
  const supportMap = {
    extensive: ['extensive', 'good'],
    professional: ['extensive', 'good'],
    documentation: ['high', 'excellent'],
    minimal: ['good', 'high'],
  };
  
  const acceptableSupport = supportMap[preference] || [];
  
  let score = 0;
  if (acceptableSupport.includes(distro.community_support)) score += 0.7;
  if (distro.professional_support && preference === 'professional') score += 0.3;
  if (acceptableSupport.includes(distro.documentation_quality)) score += 0.5;
  
  return Math.min(score, 1);
}

function scorePhilosophy(distro, preference) {
  if (!preference || preference === 'no_preference') return 1;

  const philosophyMap = {
    free_software: ['free_software'],
    freedom: ['user_freedom'],
    corporate: ['pragmatic'],
  };

  const acceptablePhilosophies = philosophyMap[preference] || [];
  return acceptablePhilosophies.includes(distro.philosophy) ? 1 : 0.6;
}

function scoreHardwareType(distro, preference) {
  if (!preference || preference === 'not_sure') return 1;

  const typeMap = {
    desktop: ['general_desktop', 'development', 'gaming', 'content_creation'],
    laptop: ['general_desktop', 'development', 'content_creation'],
    server: ['server'],
    raspberry_pi: ['server', 'general_desktop'],
    virtual_machine: ['server', 'general_desktop', 'development'],
  };

  const acceptableUseCases = typeMap[preference] || [];

  // Special handling for Raspberry Pi
  if (preference === 'raspberry_pi') {
    if (distro.id === 'raspberry-pi-os') return 1;
    if (distro.cpu_architecture.includes('arm64') || distro.cpu_architecture.includes('armhf')) return 0.9;
    return 0.1;
  }

  // Special handling for server (headless)
  if (preference === 'server') {
    if (distro.use_cases.includes('server')) return 1;
    if (distro.desktop_environments.includes('Web UI') || distro.desktop_environments.length === 0) return 0.8;
    return 0.3;
  }

  // For other types, check use case overlap
  if (distro.use_cases.some(uc => acceptableUseCases.includes(uc))) return 1;
  return 0.5;
}

function scorePrivacyLevel(distro, preference) {
  if (!preference || preference === 'not_sure' || preference === 'casual') return 1;

  // High/extreme privacy: prioritize privacy-focused distros
  if (preference === 'high' || preference === 'extreme') {
    if (distro.philosophy === 'privacy') return 1;
    if (distro.use_cases.includes('privacy')) return 0.95;
    if (distro.desktop_environments.length === 0 || distro.desktop_environments.includes('Web UI')) return 0.85;
    return 0.5;
  }

  // Enhanced privacy: some privacy features
  if (preference === 'enhanced') {
    if (distro.philosophy === 'privacy') return 1;
    if (distro.use_cases.includes('privacy')) return 0.9;
    return 0.6;
  }

  return 0.5;
}

function scoreLearningGoal(distro, preference) {
  if (!preference || preference === 'not_sure') return 1;

  // Productivity - low maintenance, stable, user-friendly
  if (preference === 'productivity') {
    if (distro.release_model === 'stable_lts' && distro.experience_level.includes('beginner')) return 1;
    if (distro.release_model === 'stable_lts') return 0.85;
    if (distro.desktop_environments.includes('GNOME') || distro.desktop_environments.includes('Cinnamon')) return 0.8;
    if (distro.community_support === 'extensive') return 0.75;
    return 0.5;
  }

  // Learning - educational, hands-on
  if (preference === 'learning') {
    if (distro.release_model === 'rolling' && !distro.experience_level.includes('beginner')) return 1;
    if (distro.release_model === 'rolling') return 0.9;
    if (distro.package_manager === 'portage' || distro.package_manager === 'pacman') return 0.85;
    if (distro.experience_level.includes('advanced') || distro.experience_level.includes('expert')) return 0.8;
    return 0.4;
  }

  // Balance - good middle ground
  if (preference === 'balance') {
    if (distro.release_model === 'semi_rolling' || distro.release_model === 'fixed_release') return 0.9;
    if (distro.experience_level.includes('intermediate')) return 0.8;
    if (distro.desktop_environments.includes('KDE Plasma') || distro.desktop_environments.includes('GNOME')) return 0.7;
    return 0.5;
  }

  return 0.5;
}

// Main scoring function
export function calculateScores(preferences) {
  return distros.map(distro => {
    let totalScore = 0;

    // Experience level (25%)
    totalScore += scoreExperienceLevel(distro, preferences.experienceLevel) * WEIGHTS.experienceLevel;

    // Use case (20%)
    totalScore += scoreUseCase(distro, preferences.useCase) * WEIGHTS.useCase;

    // Hardware (20%)
    totalScore += scoreHardware(distro, preferences.hardware) * WEIGHTS.hardware;

    // Hardware type (5%)
    totalScore += scoreHardwareType(distro, preferences.hardware?.type) * WEIGHTS.hardwareType;

    // Desktop environment (8%)
    totalScore += scoreDesktopEnvironment(distro, preferences.desktopEnvironment) * WEIGHTS.desktopEnvironment;

    // Release model (3%)
    totalScore += scoreReleaseModel(distro, preferences.releaseModel) * WEIGHTS.releaseModel;

    // Package manager (3%)
    totalScore += scorePackageManager(distro, preferences.packageManager) * WEIGHTS.packageManager;

    // Privacy level (3%)
    totalScore += scorePrivacyLevel(distro, preferences.privacyLevel) * WEIGHTS.privacyLevel;

    // Learning goal (2%)
    totalScore += scoreLearningGoal(distro, preferences.learningGoal) * WEIGHTS.learningGoal;

    // Support level (2%)
    totalScore += scoreSupportLevel(distro, preferences.supportLevel) * WEIGHTS.supportLevel;

    // Philosophy (1%)
    totalScore += scorePhilosophy(distro, preferences.philosophy) * WEIGHTS.philosophy;

    return {
      ...distro,
      score: Math.round(totalScore * 100),
    };
  });
}

// Get top recommendations with dynamic result count
export function getRecommendations(preferences) {
  const scoredDistros = calculateScores(preferences);

  // Filter out distros with very low scores (< 25%)
  const filtered = scoredDistros.filter(d => d.score >= 25);

  if (filtered.length === 0) {
    return [];
  }

  // Sort by score descending with tie-breaking
  const sorted = filtered.sort((a, b) => {
    if (b.score !== a.score) {
      return b.score - a.score;
    }
    // Tie-breaking: prefer distros with more matched attributes
    const aMatches = countMatches(a, preferences);
    const bMatches = countMatches(b, preferences);
    if (bMatches !== aMatches) {
      return bMatches - aMatches;
    }
    // Tie-breaking: prefer distros with more use cases (more versatile)
    return b.use_cases.length - a.use_cases.length;
  });

  // Determine dynamic result count based on score distribution
  const topScore = sorted[0].score;
  let resultCount;

  if (topScore >= 85) {
    resultCount = Math.min(8, sorted.length);
  } else if (topScore >= 70) {
    resultCount = Math.min(6, sorted.length);
  } else if (topScore >= 50) {
    resultCount = Math.min(5, sorted.length);
  } else {
    resultCount = Math.min(3, sorted.length);
  }

  // Apply diversity filter to avoid showing too many similar distros
  const diverseResults = applyDiversityFilter(sorted.slice(0, resultCount * 2));

  return diverseResults.slice(0, resultCount);
}

// Count how many preferences a distro matches
function countMatches(distro, preferences) {
  let matches = 0;

  // Experience level
  if (preferences.experienceLevel && distro.experience_level.includes(preferences.experienceLevel)) {
    matches++;
  }

  // Use case
  if (preferences.useCase && distro.use_cases.includes(preferences.useCase)) {
    matches++;
  }

  // Desktop environment
  if (preferences.desktopEnvironment && distro.desktop_environments.includes(preferences.desktopEnvironment)) {
    matches++;
  }

  // Release model
  if (preferences.releaseModel && distro.release_model === preferences.releaseModel) {
    matches++;
  }

  // Package manager
  if (preferences.packageManager && distro.package_manager === preferences.packageManager) {
    matches++;
  }

  // Philosophy
  if (preferences.philosophy && distro.philosophy === preferences.philosophy) {
    matches++;
  }

  return matches;
}

// Apply diversity filter to avoid showing too many similar distros
function applyDiversityFilter(distros) {
  const diverse = [];
  const seenFamilies = new Set();

  for (const distro of distros) {
    const family = distro.based_on || 'independent';

    // Allow up to 2 distros from the same family
    if (!seenFamilies.has(family) || seenFamilies.get(family) < 2) {
      diverse.push(distro);
      seenFamilies.set(family, (seenFamilies.get(family) || 0) + 1);
    }
  }

  return diverse;
}
