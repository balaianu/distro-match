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

// Sum of weights, used to normalize scores to a 0-100 scale
const WEIGHTS_TOTAL = Object.values(WEIGHTS).reduce((a, b) => a + b, 0);

// Scoring functions
export function scoreExperienceLevel(distro, preference) {
  if (!preference || preference === 'not_sure') return 1;

  const experienceMap = {
    beginner: ['beginner'],
    intermediate: ['beginner', 'intermediate'],
    advanced: ['intermediate', 'advanced'],
    expert: ['advanced', 'expert'],
  };

  const acceptableLevels = experienceMap[preference] || [];
  const exactMatch = distro.experience_level.some(level => level === preference);
  const acceptableMatch = distro.experience_level.some(level => acceptableLevels.includes(level));
  if (exactMatch) return 1;
  if (acceptableMatch) return 0.8;
  return 0.3;
}

export function scoreUseCase(distro, preference) {
  if (!preference || (Array.isArray(preference) && preference.length === 0) || preference === 'not_sure') return 1;

  const useCaseMap = {
    general_desktop: ['general_desktop'],
    development: ['general_desktop', 'development'],
    server: ['server'],
    security: ['security'],
    gaming: ['general_desktop', 'gaming'],
    content_creation: ['general_desktop', 'content_creation'],
    old_hardware: ['general_desktop', 'old_hardware'],
    privacy: ['general_desktop', 'privacy'],
  };

  const preferences = Array.isArray(preference) ? preference : [preference];
  let totalScore = 0;
  let exactMatches = 0;

  for (const pref of preferences) {
    if (pref === 'not_sure') continue;
    const acceptableUseCases = useCaseMap[pref] || [pref];
    const exactMatch = distro.use_cases.includes(pref);
    const acceptableMatch = distro.use_cases.some(uc => acceptableUseCases.includes(uc));
    if (exactMatch) {
      totalScore += 1;
      exactMatches++;
    } else if (acceptableMatch) {
      totalScore += 0.8;
    } else {
      totalScore += 0.2;
    }
  }

  const avgScore = Math.min(1, totalScore / preferences.length);
  // Bonus for high exact match ratio
  if (exactMatches === preferences.length && preferences.length > 1) {
    return Math.min(1, avgScore + 0.1);
  }
  return avgScore;
}

export function scoreHardware(distro, hardware) {
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

export function scoreDesktopEnvironment(distro, preference) {
  if (!preference || (Array.isArray(preference) && preference.length === 0) || preference === 'no_preference') return 1;

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

  const preferences = Array.isArray(preference) ? preference : [preference];
  let totalScore = 0;

  for (const pref of preferences) {
    if (pref === 'no_preference') continue;
    const acceptableDEs = deMap[pref] || [];
    if (acceptableDEs.length === 0) {
      totalScore += 0.5;
      continue;
    }
    if (distro.desktop_environments.includes('Any')) {
      totalScore += 0.9;
      continue;
    }
    const exactMatch = distro.desktop_environments.some(de => de === pref || acceptableDEs.includes(de));
    if (exactMatch) totalScore += 1;
    else totalScore += 0.3;
  }

  return Math.min(1, totalScore / preferences.length);
}

export function scoreReleaseModel(distro, preference) {
  if (!preference || preference === 'no_preference') return 1;

  const modelMap = {
    stable_lts: ['stable_lts'],
    semi_rolling: ['semi_rolling'],
    rolling: ['rolling', 'semi_rolling'],
    fixed_release: ['fixed_release', 'stable_lts'],
  };

  const acceptableModels = modelMap[preference] || [];
  const exactMatch = distro.release_model === preference;
  const acceptableMatch = acceptableModels.includes(distro.release_model);
  if (exactMatch) return 1;
  if (acceptableMatch) return 0.7;
  return 0.3;
}

export function scorePackageManager(distro, preference) {
  if (!preference || preference === 'no_preference') return 1;

  const pmMap = {
    apt: ['apt'],
    dnf: ['dnf'],
    pacman: ['pacman'],
    portage: ['portage'],
  };

  const acceptablePMs = pmMap[preference] || [];
  const exactMatch = distro.package_manager === preference;
  const acceptableMatch = acceptablePMs.includes(distro.package_manager);
  if (exactMatch) return 1;
  if (acceptableMatch) return 0.7;
  return 0.3;
}

export function scoreSupportLevel(distro, preference) {
  if (!preference || (Array.isArray(preference) && preference.length === 0) || preference === 'no_preference') return 1;

  const supportMap = {
    extensive: ['extensive', 'good'],
    professional: ['extensive', 'good'],
    documentation: ['high', 'excellent'],
    minimal: ['good', 'high'],
  };

  const preferences = Array.isArray(preference) ? preference : [preference];
  let totalScore = 0;

  for (const pref of preferences) {
    if (pref === 'no_preference') continue;
    const acceptableSupport = supportMap[pref] || [];
    let score = 0;
    if (acceptableSupport.includes(distro.community_support)) score += 0.5;
    if (distro.professional_support && pref === 'professional') score += 0.3;
    if (acceptableSupport.includes(distro.documentation_quality)) score += 0.3;
    totalScore += Math.min(score, 1);
  }

  return Math.min(1, totalScore / preferences.length);
}

export function scorePhilosophy(distro, preference) {
  if (!preference || preference === 'no_preference') return 1;

  const philosophyMap = {
    free_software: ['free_software'],
    freedom: ['user_freedom'],
    corporate: ['pragmatic'],
  };

  const acceptablePhilosophies = philosophyMap[preference] || [];
  const exactMatch = distro.philosophy === preference;
  const acceptableMatch = acceptablePhilosophies.includes(distro.philosophy);
  if (exactMatch) return 1;
  if (acceptableMatch) return 0.7;
  return 0.3;
}

export function scoreHardwareType(distro, preference) {
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

export function scorePrivacyLevel(distro, preference) {
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

export function scoreLearningGoal(distro, preference) {
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
export function calculateScores(preferences, distros) {
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
      score: Math.round((totalScore / WEIGHTS_TOTAL) * 100),
    };
  });
}

// Get top recommendations with dynamic result count
export function getRecommendations(preferences, distros) {
  const scoredDistros = calculateScores(preferences, distros);

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

  // Use case (handle arrays)
  if (preferences.useCase) {
    const ucs = Array.isArray(preferences.useCase) ? preferences.useCase : [preferences.useCase];
    if (ucs.some(uc => distro.use_cases.includes(uc))) {
      matches++;
    }
  }

  // Desktop environment (handle arrays)
  if (preferences.desktopEnvironment) {
    const des = Array.isArray(preferences.desktopEnvironment) ? preferences.desktopEnvironment : [preferences.desktopEnvironment];
    if (des.some(de => distro.desktop_environments.includes(de))) {
      matches++;
    }
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
  const familyCounts = new Map();

  for (const distro of distros) {
    const family = distro.based_on || 'independent';

    // Allow up to 2 distros from the same family
    const count = familyCounts.get(family) || 0;
    if (count < 2) {
      diverse.push(distro);
      familyCounts.set(family, count + 1);
    }
  }

  return diverse;
}
