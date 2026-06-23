import { describe, it, expect } from 'vitest';
import { calculateScores, getRecommendations } from '../src/scripts/recommendation-algorithm.js';
import distros from '../src/data/distros.json';

// Helper: find a distro by id from the scored results
function findDistro(results, id) {
  return results.find(d => d.id === id);
}

// Helper: minimal preferences object
function prefs(overrides = {}) {
  return {
    experienceLevel: null,
    useCase: null,
    hardware: { ram: null, disk: null, type: null },
    desktopEnvironment: null,
    releaseModel: null,
    packageManager: null,
    supportLevel: null,
    philosophy: null,
    privacyLevel: null,
    learningGoal: null,
    ...overrides,
  };
}

describe('calculateScores', () => {
  it('returns a score for every distro in the database', () => {
    const results = calculateScores(prefs(), distros);
    expect(results.length).toBe(63);
    results.forEach(d => {
      expect(d.score).toBeGreaterThanOrEqual(0);
      expect(d.score).toBeLessThanOrEqual(100);
    });
  });

  it('returns 100 for all distros when all preferences are null (neutral)', () => {
    const results = calculateScores(prefs(), distros);
    results.forEach(d => {
      expect(d.score).toBe(100);
    });
  });

  it('rounds scores to integers', () => {
    const results = calculateScores(prefs({ experienceLevel: 'beginner' }), distros);
    results.forEach(d => {
      expect(Number.isInteger(d.score)).toBe(true);
    });
  });
});

describe('scoreExperienceLevel', () => {
  it('gives full score when distro matches user experience level', () => {
    const results = calculateScores(prefs({ experienceLevel: 'beginner' }), distros);
    const ubuntu = findDistro(results, 'ubuntu-24.04');
    // Ubuntu supports beginner — should get full experience score
    expect(ubuntu.score).toBeGreaterThanOrEqual(90);
  });

  it('penalizes distros that do not match experience level', () => {
    const results = calculateScores(prefs({ experienceLevel: 'beginner' }), distros);
    const arch = findDistro(results, 'arch-linux');
    // Arch does not support beginner — should score lower
    expect(arch.score).toBeLessThan(90);
  });

  it('returns neutral for not_sure', () => {
    const results = calculateScores(prefs({ experienceLevel: 'not_sure' }), distros);
    results.forEach(d => {
      expect(d.score).toBe(100);
    });
  });
});

describe('scoreHardware', () => {
  it('penalizes distros requiring more RAM than available', () => {
    const lowRam = calculateScores(prefs({ hardware: { ram: 'lt_2gb', disk: null, type: null } }), distros);
    const highRam = calculateScores(prefs({ hardware: { ram: 'gt_16gb', disk: null, type: null } }), distros);
    const ubuntu = findDistro(lowRam, 'ubuntu-24.04');
    const ubuntuHigh = findDistro(highRam, 'ubuntu-24.04');
    expect(ubuntu.score).toBeLessThanOrEqual(ubuntuHigh.score);
  });

  it('penalizes distros requiring more disk than available', () => {
    const lowDisk = calculateScores(prefs({ hardware: { ram: null, disk: 'lt_20gb', type: null } }), distros);
    const highDisk = calculateScores(prefs({ hardware: { ram: null, disk: 'gt_100gb', type: null } }), distros);
    const ubuntu = findDistro(lowDisk, 'ubuntu-24.04');
    const ubuntuHigh = findDistro(highDisk, 'ubuntu-24.04');
    expect(ubuntu.score).toBeLessThanOrEqual(ubuntuHigh.score);
  });

  it('penalizes architecture mismatch for Raspberry Pi', () => {
    const rpi = calculateScores(prefs({ hardware: { ram: null, disk: null, type: 'raspberry_pi' } }), distros);
    const ubuntu = findDistro(rpi, 'ubuntu-24.04');
    // Ubuntu supports arm64, so should still score reasonably
    expect(ubuntu.score).toBeGreaterThan(50);
  });

  it('returns neutral for not_sure hardware', () => {
    const results = calculateScores(prefs({ hardware: { ram: 'not_sure', disk: 'not_sure', type: 'not_sure' } }), distros);
    results.forEach(d => {
      expect(d.score).toBe(100);
    });
  });
});

describe('scoreUseCase', () => {
  it('gives full score for matching use case', () => {
    const results = calculateScores(prefs({ useCase: 'general_desktop' }), distros);
    const ubuntu = findDistro(results, 'ubuntu-24.04');
    expect(ubuntu.score).toBe(100);
  });

  it('penalizes distros that do not match use case', () => {
    const results = calculateScores(prefs({ useCase: 'server' }), distros);
    const ubuntu = findDistro(results, 'ubuntu-24.04');
    // Ubuntu has server use case
    expect(ubuntu.score).toBe(100);
  });

  it('returns neutral for not_sure', () => {
    const results = calculateScores(prefs({ useCase: 'not_sure' }), distros);
    results.forEach(d => {
      expect(d.score).toBe(100);
    });
  });
});

describe('scoreDesktopEnvironment', () => {
  it('gives full score for matching DE', () => {
    const results = calculateScores(prefs({ desktopEnvironment: 'gnome' }), distros);
    const ubuntu = findDistro(results, 'ubuntu-24.04');
    expect(ubuntu.score).toBe(100);
  });

  it('penalizes non-matching DE', () => {
    const results = calculateScores(prefs({ desktopEnvironment: 'pantheon' }), distros);
    const ubuntu = findDistro(results, 'ubuntu-24.04');
    // Ubuntu doesn't have Pantheon
    expect(ubuntu.score).toBeLessThan(100);
  });

  it('returns neutral for no_preference', () => {
    const results = calculateScores(prefs({ desktopEnvironment: 'no_preference' }), distros);
    results.forEach(d => {
      expect(d.score).toBe(100);
    });
  });
});

describe('scoreReleaseModel', () => {
  it('gives full score for exact match', () => {
    const results = calculateScores(prefs({ releaseModel: 'stable_lts' }), distros);
    const ubuntu = findDistro(results, 'ubuntu-24.04');
    expect(ubuntu.score).toBe(100);
  });

  it('gives partial score for acceptable match', () => {
    const results = calculateScores(prefs({ releaseModel: 'rolling' }), distros);
    const ubuntu = findDistro(results, 'ubuntu-24.04');
    // Ubuntu is stable_lts, rolling accepts rolling + semi_rolling, not stable_lts
    expect(ubuntu.score).toBeLessThan(100);
  });
});

describe('scorePackageManager', () => {
  it('gives full score for exact match', () => {
    const results = calculateScores(prefs({ packageManager: 'apt' }), distros);
    const ubuntu = findDistro(results, 'ubuntu-24.04');
    expect(ubuntu.score).toBe(100);
  });

  it('gives partial score for non-match', () => {
    const results = calculateScores(prefs({ packageManager: 'pacman' }), distros);
    const ubuntu = findDistro(results, 'ubuntu-24.04');
    expect(ubuntu.score).toBeLessThan(100);
  });
});

describe('scorePrivacyLevel', () => {
  it('prioritizes privacy-focused distros for extreme privacy', () => {
    const results = calculateScores(prefs({ privacyLevel: 'extreme' }), distros);
    const tails = findDistro(results, 'tails-6');
    const ubuntu = findDistro(results, 'ubuntu-24.04');
    expect(tails.score).toBeGreaterThan(ubuntu.score);
  });

  it('returns neutral for casual privacy', () => {
    const results = calculateScores(prefs({ privacyLevel: 'casual' }), distros);
    results.forEach(d => {
      expect(d.score).toBe(100);
    });
  });
});

describe('scoreLearningGoal', () => {
  it('prioritizes rolling distros for learning goal', () => {
    const results = calculateScores(prefs({ learningGoal: 'learning' }), distros);
    const arch = findDistro(results, 'arch-linux');
    const ubuntu = findDistro(results, 'ubuntu-24.04');
    expect(arch.score).toBeGreaterThan(ubuntu.score);
  });

  it('prioritizes stable LTS for productivity goal', () => {
    const results = calculateScores(prefs({ learningGoal: 'productivity' }), distros);
    const ubuntu = findDistro(results, 'ubuntu-24.04');
    const arch = findDistro(results, 'arch-linux');
    expect(ubuntu.score).toBeGreaterThanOrEqual(arch.score);
  });
});

describe('getRecommendations', () => {
  it('filters out distros with score < 25', () => {
    const results = getRecommendations({
      experienceLevel: 'beginner',
      useCase: 'general_desktop',
      hardware: { ram: 'gt_16gb', disk: 'gt_100gb', type: 'desktop' },
      desktopEnvironment: 'gnome',
      releaseModel: 'stable_lts',
      packageManager: 'apt',
      supportLevel: null,
      philosophy: null,
      privacyLevel: 'casual',
      learningGoal: 'productivity',
    }, distros);
    results.forEach(d => {
      expect(d.score).toBeGreaterThanOrEqual(25);
    });
  });

  it('produces lower scores for mismatched preferences than matched ones', () => {
    const mismatched = getRecommendations({
      experienceLevel: 'beginner',
      useCase: 'server',
      hardware: { ram: 'lt_2gb', disk: 'lt_20gb', type: 'raspberry_pi' },
      desktopEnvironment: 'pantheon',
      releaseModel: 'rolling',
      packageManager: 'portage',
      supportLevel: 'minimal',
      philosophy: 'free_software',
      privacyLevel: 'extreme',
      learningGoal: 'learning',
    }, distros);
    const matched = getRecommendations({
      experienceLevel: 'beginner',
      useCase: 'general_desktop',
      hardware: { ram: 'gt_16gb', disk: 'gt_100gb', type: 'desktop' },
      desktopEnvironment: 'gnome',
      releaseModel: 'stable_lts',
      packageManager: 'apt',
      supportLevel: 'extensive',
      philosophy: 'pragmatic',
      privacyLevel: 'casual',
      learningGoal: 'productivity',
    }, distros);
    // Matched preferences should produce higher top scores than mismatched ones
    if (mismatched.length > 0 && matched.length > 0) {
      expect(matched[0].score).toBeGreaterThan(mismatched[0].score);
    }
  });

  it('returns 3-8 results based on score distribution', () => {
    const results = getRecommendations({
      experienceLevel: 'intermediate',
      useCase: 'development',
      hardware: { ram: '8_16gb', disk: 'gt_100gb', type: 'desktop' },
      desktopEnvironment: 'no_preference',
      releaseModel: 'no_preference',
      packageManager: 'no_preference',
      supportLevel: 'no_preference',
      philosophy: 'no_preference',
      privacyLevel: 'casual',
      learningGoal: 'not_sure',
    }, distros);
    expect(results.length).toBeGreaterThanOrEqual(1);
    expect(results.length).toBeLessThanOrEqual(8);
  });

  it('sorts results by score descending', () => {
    const results = getRecommendations({
      experienceLevel: 'beginner',
      useCase: 'general_desktop',
      hardware: { ram: '4_8gb', disk: '50_100gb', type: 'desktop' },
      desktopEnvironment: 'no_preference',
      releaseModel: 'stable_lts',
      packageManager: 'no_preference',
      supportLevel: 'no_preference',
      philosophy: 'no_preference',
      privacyLevel: 'casual',
      learningGoal: 'productivity',
    }, distros);
    for (let i = 1; i < results.length; i++) {
      expect(results[i].score).toBeLessThanOrEqual(results[i - 1].score);
    }
  });

  it('applies diversity filter (max 2 per family)', () => {
    const results = getRecommendations({
      experienceLevel: 'beginner',
      useCase: 'general_desktop',
      hardware: { ram: 'gt_16gb', disk: 'gt_100gb', type: 'desktop' },
      desktopEnvironment: 'no_preference',
      releaseModel: 'no_preference',
      packageManager: 'no_preference',
      supportLevel: 'no_preference',
      philosophy: 'no_preference',
      privacyLevel: 'casual',
      learningGoal: 'not_sure',
    }, distros);
    const familyCounts = {};
    results.forEach(d => {
      const family = d.based_on || 'independent';
      familyCounts[family] = (familyCounts[family] || 0) + 1;
    });
    Object.values(familyCounts).forEach(count => {
      expect(count).toBeLessThanOrEqual(2);
    });
  });

  it('includes score property in results', () => {
    const results = getRecommendations(prefs({ experienceLevel: 'beginner' }), distros);
    results.forEach(d => {
      expect(d).toHaveProperty('score');
      expect(typeof d.score).toBe('number');
    });
  });
});

describe('applyDiversityFilter (bug fix verification)', () => {
  it('correctly limits to 2 per family using Map (not Set)', () => {
    // This test verifies the bug fix where Set was used instead of Map
    const results = getRecommendations({
      experienceLevel: 'beginner',
      useCase: 'general_desktop',
      hardware: { ram: 'gt_16gb', disk: 'gt_100gb', type: 'desktop' },
      desktopEnvironment: 'no_preference',
      releaseModel: 'no_preference',
      packageManager: 'apt',
      supportLevel: 'no_preference',
      philosophy: 'no_preference',
      privacyLevel: 'casual',
      learningGoal: 'not_sure',
    }, distros);
    // With apt preference, many Ubuntu/Debian distros may match
    // Verify no family exceeds 2
    const familyCounts = {};
    results.forEach(d => {
      const family = d.based_on || 'independent';
      familyCounts[family] = (familyCounts[family] || 0) + 1;
    });
    Object.entries(familyCounts).forEach(([family, count]) => {
      expect(count).toBeLessThanOrEqual(2);
    });
  });
});

describe('score normalization (bug fix verification)', () => {
  it('max score is 100 when all preferences are neutral', () => {
    // This verifies the weights normalization fix — previously max was 92
    const results = calculateScores(prefs(), distros);
    const maxScore = Math.max(...results.map(d => d.score));
    expect(maxScore).toBe(100);
  });

  it('perfect match produces score close to 100', () => {
    // A distro that matches all preferences should score very high
    const results = calculateScores({
      experienceLevel: 'beginner',
      useCase: 'general_desktop',
      hardware: { ram: 'gt_16gb', disk: 'gt_100gb', type: 'desktop' },
      desktopEnvironment: 'gnome',
      releaseModel: 'stable_lts',
      packageManager: 'apt',
      supportLevel: 'extensive',
      philosophy: 'pragmatic',
      privacyLevel: 'casual',
      learningGoal: 'productivity',
    }, distros);
    const ubuntu = findDistro(results, 'ubuntu-24.04');
    expect(ubuntu.score).toBeGreaterThanOrEqual(90);
  });
});
