# Copyright 2024 Edgeless Systems GmbH
# SPDX-License-Identifier: AGPL-3.0-only

from setuptools import setup

setup(
    name="igvm-snakeoil-key",
    version="1.0.0",
    description="igvm-snakeoil-key",
    long_description=open("README.md").read(),
    long_description_content_type="text/markdown",
    packages=['gen_snakeoil_pem'],
    entry_points={
        'console_scripts': [
            'gen_snakeoil_pem = gen_snakeoil_pem.gen_snakeoil_pem:main',
        ]},
    install_requires=[
        "ecdsa",
    ]
)
