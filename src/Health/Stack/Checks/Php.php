<?php

/*
 * Copyright (c) 2026. Jakub Marszałek <marszalek.jakub@proton.me>
 * Licensed under the MIT License. See LICENSE file in the project root.
 */

declare(strict_types=1);

namespace Health\Stack\Checks;

use Health\Stack\Contracts\HealthCheckInterface;

class Php implements HealthCheckInterface
{
    private const array EXTENSIONS = [
        'pdo',
        'pdo_pgsql',
        'redis',
        'intl',
        'json'
    ];

    private const array INI_SETTINGS = [
        'memory_limit',
        'max_execution_time'
    ];

    /**
     * @var bool
     */
    private bool $allPassed = true;

    /**
     * Return php version, extensions and main ini settings.
     *
     * @return array
     */
    public function check(): array
    {
        return [
            'version' => phpversion(),
            'extensions' => $this->getLoadedExtensionsStatuses(),
            'ini_settings' => $this->getMainIniSettings()
        ];
    }

    /**
     * Check if all checks passed.
     *
     * @return bool
     */
    public function isHealthy(): bool
    {
        if (!$this->allPassed) {
            return false;
        }
        return true;
    }

    /**
     * Get loaded extension statuses
     *
     * @return array
     */
    private function getLoadedExtensionsStatuses(): array
    {
        $loadedExtensions = [];
        foreach (self::EXTENSIONS as $extension) {
            $loadedExtensions[$extension] = extension_loaded($extension);
            if (!$loadedExtensions[$extension]) {
                $this->allPassed = false;
            }
        }
        return $loadedExtensions;
    }

    /**
     * Get main ini settings
     *
     * @return array
     */
    private function getMainIniSettings(): array
    {
        $mainIniSettings = [];
        foreach (self::INI_SETTINGS as $setting) {
            $mainIniSettings[$setting] = ini_get($setting);
        }
        return $mainIniSettings;
    }
}
