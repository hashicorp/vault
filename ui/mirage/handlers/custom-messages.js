/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export default function (server) {
  server.get('/sys/config/ui/custom-messages', (schema, request) => {
    if (request.queryParams.authenticated === 'true') {
      return {
        data: {
          'key-info': {
            '01234567-89ab-cdef-0123-456789abcdef': {
              title: 'Authenticated custom message title',
              type: 'modal',
              authenticated: true,
              start_time: '2023-10-15T02:36:43.986212308Z',
              end_time: '2024-10-15T02:36:43.986212308Z',
              active: true,
            },
            '76543210-89ab-cdef-0123-456789abcdef': {
              title: 'Authenticated custom message title two',
              type: 'banner',
              authenticated: true,
              start_time: '2021-10-15T02:36:43.986212308Z',
              end_time: '2021-11-15T02:36:43.986212308Z',
              active: false,
            },
          },
          keys: ['01234567-89ab-cdef-0123-456789abcdef', '76543210-89ab-cdef-0123-456789abcdef'],
        },
      };
    }

    return {
      data: {
        'key-info': {
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
      data: {
        'key-info': {
          '01234567-89ab-cdef-0123-456789abcdef': {
            title: 'Unauthenticated Title One',
            message:
              'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur nulla augue, placerat quis risus blandit, molestie imperdiet massa. Sed blandit rutrum odio quis varius. Fusce purus orci, maximus ac libero.',
            type: 'modal',
            authenticated: false,
            start_time: '2023-10-15T02:36:43.986212308Z',
            end_time: '2024-10-15T02:36:43.986212308Z',
            options: {},
          },
          '76543210-89ab-cdef-0123-456789abcdef': {
            title: 'Unauthenticated Title Two',
            message:
              'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur nulla augue, placerat quis risus blandit, molestie imperdiet massa. Sed blandit rutrum odio quis varius. Fusce purus orci, maximus ac libero.',
            type: 'banner',
            authenticated: false,
            start_time: '2021-10-15T02:36:43.986212308Z',
            end_time: '2031-10-15T02:36:43.986212308Z',
            options: {},
          },
        },
        keys: ['01234567-89ab-cdef-0123-456789abcdef', '76543210-89ab-cdef-0123-456789abcdef'],
      },
    };
  });

  server.get('/sys/internal/ui/authenticated-messages', () => {
    return {
      data: {
        'key-info': {
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
