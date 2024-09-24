/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export default function (server) {
  server.get('/sys/config/ui/custom-messages', (schema, request) => {
    if (request.queryParams.authenticated && JSON.parse(request.queryParams.authenticated)) {
      return {
        data: {
          key_info: {
            '01234567-89ab-cdef-0123-456789abcdef': {
              title: 'Has expiration date',
              type: 'modal',
              authenticated: true,
              start_time: '2023-10-15T02:36:43.986212308Z',
              end_time: '2023-12-17T02:36:43.986212308Z',
              active: true,
            },
            '22234567-89ab-cdef-0123-456789abcdef': {
              title: 'No expiration date',
              type: 'modal',
              authenticated: true,
              start_time: '2023-10-15T02:36:43.986212308Z',
              end_time: '',
              active: true,
            },
            '76543210-89ab-cdef-0123-456789abcdef': {
              title: 'Inactive message',
              type: 'banner',
              authenticated: true,
              start_time: '2023-10-15T02:36:43.986212308Z',
              end_time: '2023-11-15T02:36:43.986212308Z',
              active: false,
            },
            '11543210-89ab-cdef-0123-456789abcdef': {
              title: 'Inactive, but start time is past current date',
              type: 'banner',
              authenticated: true,
              start_time: '2024-10-15T02:36:43.986212308Z',
              end_time: '2024-11-15T02:36:43.986212308Z',
              active: false,
            },
          },
          keys: [
            '01234567-89ab-cdef-0123-456789abcdef',
            '22234567-89ab-cdef-0123-456789abcdef',
            '76543210-89ab-cdef-0123-456789abcdef',
            '11543210-89ab-cdef-0123-456789abcdef',
          ],
        },
      };
    }

    return {
      data: {
        key_info: {
          '8d6ba39-5c23-50af-3d79-76c26a2845f49': {
            title: 'Unauthenticated custom message title',
            type: 'modal',
            authenticated: false,
            start_time: '2023-10-15T02:36:43.986212308Z',
            end_time: '2024-10-15T02:36:43.986212308Z',
            active: true,
          },
          '281e580-da16-89c5-4666-16480e4b7c11d': {
            title: 'Unauthenticated custom message title two',
            type: 'banner',
            authenticated: false,
            start_time: '2021-10-15T02:36:43.986212308Z',
            end_time: '2021-11-15T02:36:43.986212308Z',
            active: false,
          },
        },
        keys: ['8d6ba39-5c23-50af-3d79-76c26a2845f49', '281e580-da16-89c5-4666-16480e4b7c11d'],
      },
    };
  });

  server.post('/sys/config/ui/custom-messages', () => {
    return {
      id: '01234567-89ab-cdef-0123-456789abcdef',
      data: {
        active: true,
        start_time: '2023-10-15T02:36:43.986212308Z',
        end_time: '2024-10-15T02:36:43.986212308Z',
        type: 'modal',
        authenticated: false,
      },
    };
  });

  server.get('/sys/internal/ui/unauthenticated-messages', () => {
    return {
      request_id: '664fbad0-fcd8-9023-4c5b-81a7962e9f4b',
      lease_id: '',
      renewable: false,
      lease_duration: 0,
      data: {
        key_info: {
          '02180e3f-bd5b-a851-bcc9-6f7983806df0': {
            authenticated: false,
            end_time: null,
            link: {
              title: '',
            },
            message: 'aGVsbG8gd29ybGQgaGVsbG8gd29scmQ=',
            options: null,
            start_time: '2024-01-04T08:00:00Z',
            title: 'Banner title',
            type: 'banner',
          },
          'a7d7d9b1-a1ca-800c-17c5-0783be88e29c': {
            authenticated: false,
            end_time: null,
            link: {
              title: '',
            },
            message: 'aGVyZSBpcyBhIGNvb2wgbWVzc2FnZQ==',
            options: null,
            start_time: '2024-01-01T08:00:00Z',
            title: 'Modal title',
            type: 'modal',
          },
        },
        keys: ['02180e3f-bd5b-a851-bcc9-6f7983806df0', 'a7d7d9b1-a1ca-800c-17c5-0783be88e29c'],
      },
      wrap_info: null,
      warnings: null,
      auth: null,
      mount_type: '',
    };
  });

  server.get('/sys/internal/ui/authenticated-messages', () => {
    return {
      data: {
        key_info: {
          '6543210-89ab-cdef-0123-456780abcieh': {
            title: 'Authenticated Title One',
            message:
              'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur nulla augue, placerat quis risus blandit, molestie imperdiet massa. Sed blandit rutrum odio quis varius. Fusce purus orci, maximus ac libero.',
            type: 'modal',
            authenticated: true,
            start_time: '2023-10-15T02:36:43.986212308Z',
            end_time: '2024-10-15T02:36:43.986212308Z',
            options: {},
          },
          '00123858-89ab-cdef-0123-789037ejhdgt': {
            title: 'Authenticated Title One',
            message:
              'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur nulla augue, placerat quis risus blandit, molestie imperdiet massa. Sed blandit rutrum odio quis varius. Fusce purus orci, maximus ac libero.',
            type: 'banner',
            authenticated: true,
            start_time: '2021-10-15T02:36:43.986212308Z',
            end_time: '2031-10-15T02:36:43.986212308Z',
            options: {},
          },
        },
        keys: ['6543210-89ab-cdef-0123-456780abcieh', '00123858-89ab-cdef-0123-789037ejhdgt'],
      },
    };
  });
}
