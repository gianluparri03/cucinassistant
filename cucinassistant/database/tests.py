from cucinassistant.exceptions import CAError
from cucinassistant.database import *

from unittest import TestCase


class TestDatabase(TestCase):
    @classmethod
    def setUpClass(cls):
        cls.db = init_db(testing=True)
        print('- Database initialized')

        cls.user1 = create_user('user1', 'user1@email.com', 'password1')
        cls.user2 = create_user('user2', 'user2@email.com', 'password2')
        print('- Test users created')

        print('Starting tests...\n')

    def setUp(self):
        self.db = TestDatabase.db
        self.user1 = TestDatabase.user1
        self.user2 = TestDatabase.user2

    def test_users(self):
        # Tests for create_user
        self.assertRaisesRegex(CAError, 'Nome utente non valido \(lunghezza minima', create_user, '', '', '')
        self.assertRaisesRegex(CAError, 'Nome utente non valido \(solo lettere', create_user, 'antonio&', '', '')
        self.assertRaisesRegex(CAError, 'Password non valida \(lunghezza minima', create_user, 'antonio', '', '')
        antonio = create_user('antonio', 'antonio@email.com', 'passwordA')
        self.assertRaisesRegex(CAError, 'Email non disponibile', create_user, 'antonello', 'antonio@email.com', 'passwordA')
        self.assertRaisesRegex(CAError, 'Nome utente non disponibile', create_user, 'antonio', 'antonello@email.com', 'passwordA')
        antonello = create_user('antonello', 'antonello@email.com', 'passwordA')
        fake_user = 0

        # Tests for login_user
        self.assertEqual(antonio, login_user('antonio', 'passwordA'))
        self.assertRaisesRegex(CAError, 'Credenziali non valide', login_user, 'antonio', 'passwordB')
        self.assertRaisesRegex(CAError, 'Credenziali non valide', login_user, 'antonino', 'passwordB')

        # Tests for get_user_data
        antonio_data = get_user_data(antonio)
        self.assertEqual(antonio_data.uid, antonio)
        self.assertEqual(antonio_data.username, 'antonio')
        self.assertEqual(antonio_data.email, 'antonio@email.com')
        self.assertEqual(antonio_data, get_user_data(0, email='antonio@email.com'))
        self.assertRaisesRegex(CAError, 'Utente sconosciuto', get_user_data, fake_user)

        # Tests for change_user_username
        self.assertRaisesRegex(CAError, 'Utente sconosciuto', change_user_username, fake_user, 'nuovo_username')
        change_user_username(antonello, 'antonello')
        self.assertEqual(get_user_data(antonello).username, 'antonello')
        self.assertRaisesRegex(CAError, 'Nome utente non disponibile', change_user_username, antonello, 'antonio')
        change_user_username(antonello, 'antonellino')
        self.assertEqual(get_user_data(antonello).username, 'antonellino')

        # Tests for change_user_email
        self.assertRaisesRegex(CAError, 'Utente sconosciuto', change_user_email, fake_user, 'fake_email')
        change_user_email(antonello, 'antonello@email.com')
        self.assertEqual(get_user_data(antonello).email, 'antonello@email.com')
        self.assertRaisesRegex(CAError, 'Email non disponibile', change_user_email, antonello, 'antonio@email.com')
        change_user_email(antonello, 'antonellino@email.com')
        self.assertEqual(get_user_data(antonello).email, 'antonellino@email.com')

        # Tests for change_user_password
        self.assertRaisesRegex(CAError, 'Utente sconosciuto', change_user_password, fake_user, 'fake_password', 'fake_password2')
        self.assertRaisesRegex(CAError, 'Credenziali non valide', change_user_password, antonello, 'vecchia', 'nuova')
        change_user_password(antonello, 'passwordA', 'passwordB')
        self.assertTrue(ph.verify(get_user_data(antonello).password, 'passwordB'))

        # Tests for generate_user_token
        token = generate_user_token(antonello)
        self.assertTrue(ph.verify(get_user_data(antonello).token, token))
        self.assertRaisesRegex(CAError, 'Utente sconosciuto', generate_user_token, fake_user)

        # Test for reset_user_password
        self.assertRaisesRegex(CAError, 'Utente sconosciuto', reset_user_password, 'fake@email.com', 'fake_token', 'fake_password')
        self.assertRaisesRegex(CAError, 'Errore durante la reimpostazione', reset_user_password, 'antonellino@email.com', 'token', 'new_password')
        reset_user_password('antonellino@email.com', token, 'passwordA')
        self.assertTrue(ph.verify(get_user_data(antonello).password, 'passwordA'))
        self.assertIsNone(get_user_data(antonello).token)

        # Tests for delete_user
        self.assertRaisesRegex(CAError, 'Utente sconosciuto', delete_user, fake_user, 'token')
        self.assertRaisesRegex(CAError, 'Errore durante la cancellazione', delete_user, antonello, 'token')
        self.assertRaisesRegex(CAError, 'Errore durante la cancellazione', delete_user, antonello, token)
        token = generate_user_token(antonello)
        delete_user(antonello, token)
        self.assertRaisesRegex(CAError, 'Credenziali non valide', login_user, 'antonello', 'passwordA')

        # Tests for get_users_number
        self.assertEqual(get_users_number(), 3)

        # Tests for get_users_email
        self.assertEqual(set(get_users_emails()), {'user1@email.com', 'user2@email.com', 'antonio@email.com'})

    def test_menus(self):
        # Tests for get_user_menu
        self.assertRaisesRegex(CAError, 'Utente sconosciuto', get_user_menu, 0)
        self.assertEqual(get_user_menu(self.user1), Menu(menu=[]*14, prev=None, next=None))
