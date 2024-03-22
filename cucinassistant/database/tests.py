from cucinassistant.exceptions import CAError
from cucinassistant.database import init_db, create_user
from cucinassistant.database.users_test import TestUsers
from cucinassistant.database.menus_test import TestMenus
from cucinassistant.database.lists_test import TestLists

from unittest import TestCase


class TestDatabase(TestCase):
    @classmethod
    def setUpClass(cls):
        init_db(testing=True)
        print('- Database initialized')

        cls.francesco = create_user('francesco', 'francesco@email.com', 'password1')
        cls.giovanna = create_user('giovanna', 'giovanna@email.com', 'password2')
        print('- Test users created')

        print('Starting tests...\n')

    def setUp(self):
        self.francesco = TestDatabase.francesco
        self.giovanna = TestDatabase.giovanna
        self.fake_user = 0

    def test_users(self):
        # See cucinassistant/database/users_test.py
        TestUsers(self, self.fake_user)

    def test_menus(self):
        # See cucinassistant/database/menus_test.py
        TestMenus(self, self.francesco, self.giovanna, self.fake_user)

    def test_lists(self):
        # See cucinassistant/database/lists_test.py
        TestLists(self, self.francesco, self.giovanna, self.fake_user)
