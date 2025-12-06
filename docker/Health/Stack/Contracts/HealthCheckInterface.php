<?php

/*
 * Copyright (c) 2025. Jakub Marszałek / hesoyamcc
 * All Rights Reserved.
 */

declare(strict_types=1);

namespace Health\Stack\Contracts;

interface HealthCheckInterface
{
    /**
     * Perform the health check
     *
     * @return array
     */
    public function check(): array;

    /**
     * Check a service's health status
     *
     * @return bool
     */
    public function isHealthy(): bool;
}
