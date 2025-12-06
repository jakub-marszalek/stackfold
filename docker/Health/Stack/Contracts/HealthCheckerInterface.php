<?php

/*
 * Copyright (c) 2025. Jakub Marszałek / hesoyamcc
 * All Rights Reserved.
 */

declare(strict_types=1);

namespace Health\Stack\Contracts;

interface HealthCheckerInterface
{
    /**
     * Register a health check
     *
     * @param string $name
     * @param HealthCheckInterface $check
     * @return void
     */
    public function register(string $name, HealthCheckInterface $check): void;

    /**
     * Run all registered health checks
     *
     * @return array
     */
    public function run(): array;

    /**
     * Summarize results
     *
     * @return void
     */
    public function summarize(): void;
}
