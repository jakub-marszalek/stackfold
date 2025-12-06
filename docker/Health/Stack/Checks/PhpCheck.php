<?php

/*
 * Copyright (c) 2025. Jakub Marszałek / hesoyamcc
 * All Rights Reserved.
 */

declare(strict_types=1);

namespace Health\Stack\Checks;

use Health\Stack\Contracts\HealthCheckInterface;

class PhpCheck implements HealthCheckInterface
{
    /**
     * List of extensions to check
     *
     * @var array|string[]
     */
    private array $extensions = [
        'pdo',
        'pdo_pgsql',
        'redis',
        'intl',
        'json'
    ];

    /**
     * List of main ini settings to check
     *
     * @var array|string[]
     */
    private array $iniSettings = [
        'memory_limit',
        'max_execution_time'
    ];

    /**
     * @var bool
     */
    private bool $allPassed = true;

    /**
     * Return php version, extensions and main ini settings
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
     * Check if all checks passed
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
        foreach ($this->extensions as $extension) {
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
        foreach ($this->iniSettings as $setting) {
            $mainIniSettings[$setting] = ini_get($setting);
        }
        return $mainIniSettings;
    }
}
