<?php

declare(strict_types=1);

namespace Health\Stack\Contracts;

interface HealthCheckInterface
{
    /**
     * Execute the health check.
     *
     * @return array
     */
    public function check(): array;

    /**
     * Check the health status for given service.
     *
     * @return bool
     */
    public function isHealthy(): bool;
}
