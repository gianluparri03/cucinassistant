from cucinassistant.database import init_db, create_user
from cucinassistant.exceptions import CAError

from unittest import TestCase


class SubTest:
    def __init__(self, tester, francesco, giovanna, fake_user):
        self.t = tester
        self.francesco = francesco
        self.giovanna = giovanna
        self.fake_user = fake_user

        # Executes all the subtests IN ORDER
        for name in sorted(dir(self)):
            if name[0] == 'S':
                with self.t.subTest(name):
                    getattr(self, name)()


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
        from cucinassistant.database.tests.users import TestUsers
        TestUsers(self, self.fake_user)

    def test_menus(self):
        from cucinassistant.database.tests.menus import TestMenus
        TestMenus(self, self.francesco, self.giovanna, self.fake_user)

    def test_shopping(self):
        from cucinassistant.database.tests.shopping import TestShopping
        TestShopping(self, self.francesco, self.giovanna, self.fake_user)

    def test_storage(self):
        from cucinassistant.database.tests.storage import TestStorage
        TestStorage(self, self.francesco, self.giovanna, self.fake_user)
