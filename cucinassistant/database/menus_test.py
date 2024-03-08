from cucinassistant.exceptions import CAError
from cucinassistant.database import *


class TestMenus:
    def __init__(self, tester, francesco, giovanna, fake_user):
        self.t = tester
        self.francesco = francesco
        self.giovanna = giovanna
        self.fake_user = fake_user
        self.fake_menu = 0

        # Executes all the subtests IN ORDER
        for name in sorted(dir(self)):
            if name[0] == 'S':
                with self.t.subTest(sub=int(name[1:3])):
                    getattr(self, name)()


    def S00_create_menu(self):
        # Tests for create_menu
        create_menu(self.francesco)
        create_menu(self.francesco)
        create_menu(self.francesco)
        self.t.assertEqual(get_menu(self.francesco, 1).next, 2)
        self.t.assertEqual(get_menu(self.francesco, 2).prev, 1)
        self.t.assertEqual(get_menu(self.francesco, 2).next, 3)
        self.t.assertEqual(get_menu(self.francesco, 3).prev, 2)
        self.t.assertEqual(get_menu(self.francesco, 3).menu.count(';'), 13)

    def S01_get_menu(self):
        # Tests for get_menu
        self.t.assertRaisesRegex(CAError, 'Utente sconosciuto', get_menu, self.fake_user)
        self.t.assertRaisesRegex(CAError, 'Menu non trovato', get_menu, self.giovanna, 7)
        self.t.assertRaisesRegex(CAError, 'Menu non trovato', get_menu, self.giovanna, 0)
        self.t.assertEqual(get_menu(self.francesco).mid, 3)
        self.t.assertEqual(get_menu(self.giovanna).mid, 0)

    def S02_update_menu(self):
        # Tests for update_menu
        self.t.assertRaisesRegex(CAError, 'Utente sconosciuto', update_menu, self.fake_user, '')
        self.t.assertRaisesRegex(CAError, 'Menu non trovato', update_menu, self.giovanna, 7, '')
        self.t.assertRaisesRegex(CAError, 'Menu non trovato', update_menu, self.giovanna, 0, '')
        self.t.assertRaisesRegex(CAError, 'Menu non valido', update_menu, self.francesco, 1, 'pippo')
        update_menu(self.francesco, 1, '0;1;2;3;4;5;6;7;8;9;10;11;12;13')
        self.t.assertEqual(get_menu(self.francesco, 1).menu, '0;1;2;3;4;5;6;7;8;9;10;11;12;13')
